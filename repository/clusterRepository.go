package repository

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const currentClusterFileName = "CURRENT"

type ClusterRepository interface {
	Current() (*string, error)
	SetCurrent(name string) error
	Load(name string) (*Cluster, error)
	LoadAll() ([]*Cluster, error)
	Save(cluster *Cluster) error
	Remove(name string) error
	Exists(name string) bool
}

type clusterRepository struct {
	path string
}

func New(path string) (ClusterRepository, error) {
	if err := checkDirExists(path); err != nil {
		return nil, err
	}

	return &clusterRepository{
		path: path,
	}, nil
}

func (r *clusterRepository) LoadAll() ([]*Cluster, error) {
	var result []*Cluster

	files, err := ioutil.ReadDir(r.path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read cluster repository dir '%s'", r.path)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		fileName := file.Name()
		cluster, err := r.Load(fileName)
		if err != nil {
			glog.Warningf("failed to load cluster '%s' in repository '%s'", fileName, r.path)
			continue
		}

		result = append(result, cluster)
	}

	return result, nil
}

func (r *clusterRepository) Load(name string) (*Cluster, error) {
	clusterPath := path.Join(r.path, name)
	return newCluster(clusterPath)
}

func (r *clusterRepository) Current() (*string, error) {
	filePath := path.Join(r.path, currentClusterFileName)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to load current cluster")
	}

	clusterName := string(bytes)
	return &clusterName, nil
}

func (r *clusterRepository) SetCurrent(name string) error {
	filePath := path.Join(r.path, currentClusterFileName)
	data := []byte(name)

	err := ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to set cluster '%s' as current cluster", name)
	}

	return nil
}

func (r *clusterRepository) Save(cluster *Cluster) error {
	clusterName := cluster.Name()
	if clusterName == "" {
		return errors.New("invalid cluster name")
	}

	clusterPath := path.Join(r.path, clusterName)
	return cluster.Save(clusterPath)
}

func (r *clusterRepository) Remove(name string) error {
	if name == "" {
		return errors.New("invalid cluster name")
	}

	clusterPath := path.Join(r.path, name)
	err := os.RemoveAll(clusterPath)
	if err != nil {
		return errors.Wrapf(err, "failed to remove cluster '%s'", name)
	}

	return nil
}

func (r *clusterRepository) Exists(name string) bool {
	clusterPath := path.Join(r.path, name)
	return checkDirExists(clusterPath) == nil
}

func checkDirExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errors.Wrapf(err, "invalid cluster repository path '%s'", path)
	}

	if !info.IsDir() {
		return errors.Errorf("invalid cluster repository path '%s' - not a directory", path)
	}

	return nil
}
