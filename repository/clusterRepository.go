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
	Current() (*Cluster, error)
	SetCurrent(name string) error
	Load(name string) (*Cluster, error)
	LoadAll() ([]*Cluster, error)
	Save(cluster Cluster) error
	Remove(name string) error
	Exists(name string) bool
}

type clusterRepository struct {
	path           string
	currentCluster string //TODO
}

func New(path string) (ClusterRepository, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, errors.Wrapf(err, "repository: failed to initialize cluster repository '%s'", path)
	}

	return &clusterRepository{
		path: path,
	}, nil
}

func (r *clusterRepository) LoadAll() ([]*Cluster, error) {
	var result []*Cluster

	files, err := ioutil.ReadDir(r.path)
	if err != nil {
		return nil, errors.Wrapf(err, "repository: failed to read cluster repository dir '%s'", r.path)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		fileName := file.Name()
		cluster, err := r.Load(fileName)
		if err != nil {
			glog.Warningf("repository: failed to load cluster '%s' in repository '%s'", fileName, r.path)
			continue
		}

		result = append(result, cluster)
	}

	return result, nil
}

func (r *clusterRepository) Load(name string) (*Cluster, error) {
	clusterPath := path.Join(r.path, name)
	return loadCluster(clusterPath)
}

func (r *clusterRepository) Current() (*Cluster, error) {
	filePath := path.Join(r.path, currentClusterFileName)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "repository: failed to load current cluster")
	}

	return r.Load(string(bytes))
}

func (r *clusterRepository) SetCurrent(name string) error {
	filePath := path.Join(r.path, currentClusterFileName)
	data := []byte(name)

	err := ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return errors.Wrapf(err, "repository: failed to set cluster '%s' as current cluster", name)
	}

	return nil
}

func (r *clusterRepository) Save(cluster Cluster) error {
	if err := cluster.Validate(); err != nil {
		return errors.Wrap(err, "repository: failed to save cluster")
	}

	clusterPath := path.Join(r.path, cluster.Name)
	if err := os.MkdirAll(clusterPath, 0755); err != nil {
		return errors.Wrapf(err, "repository: failed to create cluster directory '%s'", clusterPath)
	}

	return cluster.save(clusterPath)
}

func (r *clusterRepository) Remove(name string) error {
	if name == "" {
		return errors.New("repository: invalid cluster name")
	}

	clusterPath := path.Join(r.path, name)
	err := os.RemoveAll(clusterPath)
	if err != nil {
		return errors.Wrapf(err, "repository: failed to remove cluster '%s'", name)
	}

	return nil
}

func (r *clusterRepository) Exists(name string) bool {
	clusterPath := path.Join(r.path, name)

	_, err := os.Stat(clusterPath)
	if err != nil {
		return !os.IsNotExist(err)
	}

	return true
}
