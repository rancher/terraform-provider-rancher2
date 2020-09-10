package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func resourceRancher2CatalogV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2CatalogV2Create,
		Read:   resourceRancher2CatalogV2Read,
		Update: resourceRancher2CatalogV2Update,
		Delete: resourceRancher2CatalogV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2CatalogV2Import,
		},
		Schema: catalogV2Fields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2CatalogV2Create(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	catalog := expandCatalogV2(d)

	log.Printf("[INFO] Creating Catalog V2 %s at %s", name, clusterID)

	client, err := meta.(*Config).catalogV2Client(clusterID)
	if err != nil {
		return err
	}
	obj, err := client.Create("", catalog)
	if err != nil {
		return err
	}
	newCatalog := obj.(*v1.ClusterRepo)
	id := string(newCatalog.ObjectMeta.UID)
	d.SetId(id)

	timeout := int64(d.Timeout(schema.TimeoutCreate).Seconds())
	listOption := metaV1.ListOptions{
		TypeMeta:        newCatalog.TypeMeta,
		Watch:           true,
		ResourceVersion: newCatalog.ObjectMeta.ResourceVersion,
		TimeoutSeconds:  &timeout,
	}
	watcher, err := client.Watch("", listOption)
	for {
		select {
		case event, open := <-watcher.ResultChan():
			if open {
				if event.Type == watch.Added || event.Type == watch.Modified {
					if repo, ok := event.Object.DeepCopyObject().(*v1.ClusterRepo); ok && len(repo.Status.URL) > 0 && len(repo.Status.IndexConfigMapName) > 0 {
						watcher.Stop()
						return resourceRancher2CatalogV2Read(d, meta)
					}
				}
				continue
			}
			return fmt.Errorf("[ERROR] waiting for catalog V2 (%s) to be added", name)
		}
	}
}

func resourceRancher2CatalogV2Read(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Catalog V2 %s at %s", name, clusterID)

	client, err := meta.(*Config).catalogV2Client(clusterID)
	if err != nil {
		return err
	}
	obj, err := client.Get(name, "", metaV1.GetOptions{ResourceVersion: d.Get("resource_version").(string)})
	if err != nil {
		if errors.IsNotFound(err) || errors.IsForbidden(err) {
			log.Printf("[INFO] Catalog V2 %s not found at %s", name, clusterID)
			d.SetId("")
			return nil
		}
	}
	return flattenCatalogV2(d, obj.(*v1.ClusterRepo))
}

func resourceRancher2CatalogV2Update(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	catalog := expandCatalogV2(d)
	log.Printf("[INFO] Updating Catalog V2 %s at %s", name, clusterID)

	client, err := meta.(*Config).catalogV2Client(clusterID)
	if err != nil {
		return err
	}
	obj, err := client.Update("", catalog)
	if err != nil {
		return err
	}
	newCatalog := obj.(*v1.ClusterRepo)
	timeout := int64(d.Timeout(schema.TimeoutUpdate).Seconds())
	listOption := metaV1.ListOptions{
		TypeMeta:        newCatalog.TypeMeta,
		Watch:           true,
		ResourceVersion: newCatalog.ObjectMeta.ResourceVersion,
		TimeoutSeconds:  &timeout,
	}
	watcher, err := client.Watch("", listOption)
	for {
		select {
		case event, open := <-watcher.ResultChan():
			if open {
				if event.Type == watch.Added || event.Type == watch.Modified {
					if repo, ok := event.Object.DeepCopyObject().(*v1.ClusterRepo); ok && len(repo.Status.URL) > 0 && len(repo.Status.IndexConfigMapName) > 0 {
						watcher.Stop()
						return resourceRancher2CatalogV2Read(d, meta)
					}
				}
				continue
			}
			return fmt.Errorf("[ERROR] waiting for catalog V2 (%s) to be updated", name)
		}
	}
}

func resourceRancher2CatalogV2Delete(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting Catalog V2 %s at %s", name, clusterID)

	client, err := meta.(*Config).catalogV2Client(clusterID)
	if err != nil {
		return err
	}
	obj, err := client.Get(name, "", metaV1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) || errors.IsForbidden(err) {
			log.Printf("[INFO] Catalog V2 %s not found at %s", name, clusterID)
			d.SetId("")
			return nil
		}
		return err
	}
	catalog := obj.(*v1.ClusterRepo)
	err = client.Delete(name, "", nil)
	if err != nil {
		return fmt.Errorf("Error removing Catalog V2 %s: %s", name, err)
	}

	timeout := int64(d.Timeout(schema.TimeoutDelete).Seconds())
	listOption := metaV1.ListOptions{
		TypeMeta:        catalog.TypeMeta,
		Watch:           true,
		ResourceVersion: catalog.ObjectMeta.ResourceVersion,
		TimeoutSeconds:  &timeout,
	}
	watcher, err := client.Watch("", listOption)
	for {
		select {
		case event, open := <-watcher.ResultChan():
			if open {
				if event.Type == watch.Deleted {
					d.SetId("")
					watcher.Stop()
					return nil
				}
				continue
			}
			return fmt.Errorf("[ERROR] waiting for catalog V2 (%s) to be deleted", name)
		}
	}
}
