package repository

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/bcrusu/kcm/util"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const currentClusterFileName = "CURRENT"
const clusterFileExtension = ".json"

type ClusterRepository interface {
	Current() (*Cluster, error)
	SetCurrent(name string) error
	Load(name string) (*Cluster, error)
	LoadAll() ([]*Cluster, error)
	Save(cluster Cluster) error
	Remove(name string) error
	Exists(name string) (bool, error)
}

type clusterRepository struct {
	path           string
	currentCluster *string
}

func New(path string) (ClusterRepository, error) {
	if err := util.CreateDirectoryPath(path); err != nil {
		return nil, errors.Wrapf(err, "repository: failed to initialize cluster repository '%s'", path)
	}

	result := &clusterRepository{
		path: path,
	}

	if err := result.loadCurrentClusterName(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *clusterRepository) LoadAll() ([]*Cluster, error) {
	var result []*Cluster

	files, err := ioutil.ReadDir(r.path)
	if err != nil {
		return nil, errors.Wrapf(err, "repository: failed to read cluster repository dir '%s'", r.path)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if !strings.HasSuffix(fileName, clusterFileExtension) {
			continue
		}

		cluster, err := loadCluster(path.Join(r.path, fileName))
		if err != nil {
			glog.Warningf("repository: failed to load cluster from file '%s'", fileName)
			continue
		}

		result = append(result, cluster)
	}

	return result, nil
}

func (r *clusterRepository) Load(name string) (*Cluster, error) {
	filePath := r.clusterFile(name)

	exists, err := util.FileExists(filePath)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, nil
	}

	return loadCluster(filePath)
}

func (r *clusterRepository) Current() (*Cluster, error) {
	if r.currentCluster == nil {
		return nil, nil
	}

	return r.Load(*r.currentCluster)
}

func (r *clusterRepository) SetCurrent(name string) error {
	if name == "" {
		return r.clearCurrentClusterName()
	}

	filePath := path.Join(r.path, currentClusterFileName)
	data := []byte(name)

	err := util.WriteFile(filePath, data)
	if err != nil {
		return errors.Wrapf(err, "repository: failed to set cluster '%s' as current cluster", name)
	}

	r.currentCluster = &name
	return nil
}

func (r *clusterRepository) Save(cluster Cluster) error {
	if err := cluster.Validate(); err != nil {
		return errors.Wrap(err, "repository: failed to save cluster")
	}

	filePath := r.clusterFile(cluster.Name)
	if err := cluster.save(filePath); err != nil {
		return err
	}

	if r.currentCluster == nil {
		return r.SetCurrent(cluster.Name)
	}

	return nil
}

func (r *clusterRepository) Remove(name string) error {
	if name == "" {
		return errors.New("repository: invalid cluster name")
	}

	filePath := r.clusterFile(name)
	err := os.Remove(filePath)
	if err != nil {
		return errors.Wrapf(err, "repository: failed to remove cluster '%s'", name)
	}

	if name == *r.currentCluster {
		return r.clearCurrentClusterName()
	}

	return nil
}

func (r *clusterRepository) Exists(name string) (bool, error) {
	filePath := r.clusterFile(name)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r *clusterRepository) loadCurrentClusterName() error {
	filePath := path.Join(r.path, currentClusterFileName)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return errors.Wrap(err, "repository: failed to load current cluster")
	}

	name := strings.TrimSpace(string(bytes))
	r.currentCluster = &name
	return nil
}

func (r *clusterRepository) clearCurrentClusterName() error {
	filePath := path.Join(r.path, currentClusterFileName)
	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return errors.Wrap(err, "repository: failed to clear current cluster name")
	}

	r.currentCluster = nil
	return nil
}

func (r *clusterRepository) clusterFile(clusterName string) string {
	fileName := clusterName + clusterFileExtension
	return path.Join(r.path, fileName)
}
