package service

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go-openstackclient/configs"
	"go-openstackclient/consts"
    "go-openstackclient/internal"
    "go-openstackclient/internal/entity"
	"log"
)

type Glance struct {
    internal.Request
	projectId          string
	headers            map[string]string
	tag                string
}

func NewGlance(token, projectId string, client *fasthttp.Client) *Glance {
	return &Glance{
		Request: internal.Request{
			UrlPrefix: fmt.Sprintf("http://%s:%d/v2", configs.CONF.Host, consts.GlancePort),
			Client: client,
		},
		projectId: projectId,
		headers: map[string]string{"X-Auth-Token": token},
		tag: configs.CONF.Host + "_",
	}
}

func (g *Glance) CreateImage(reqBody string) string {
	createSuffix := "/images"
	resp := g.Post(g.headers, createSuffix, reqBody)

	var image entity.ImageMap
	_ = json.Unmarshal(resp, &image)
	//cache.RedisClient.SetMap(g.tag + consts.Images, image.Id, image)
	log.Println("==============Create image success", image.Id)
	return image.Id
}

func (g *Glance) SetImageProtected(imageId string, protected bool) string {
	// property-->[private, public, shared, community]
	createSuffix := fmt.Sprintf("/images/%s", imageId)
	formatter := `[{"path": "/protected", "value": %+v, "op": "replace"}]`
	reqBody := fmt.Sprintf(formatter, protected)
	resp := g.Patch(g.headers, createSuffix, reqBody)
	var image entity.ImageMap
	_ = json.Unmarshal(resp, &image)
	//cache.RedisClient.SetMap(g.tag + consts.Images, image.Id, image)
	log.Printf("==============Set image %s protected %v success", image.Id, protected)
	return image.Id
}

func (g *Glance) SetImageVisibilityProperty(imageId string, property string) string {
	// property-->[private, public, shared, community]
	createSuffix := fmt.Sprintf("/images/%s", imageId)
	formatter := `[{"path": "/visibility", "value": "%+v", "op": "replace"}]`
	reqBody := fmt.Sprintf(formatter, property)
	resp := g.Patch(g.headers, createSuffix, reqBody)
	var image entity.ImageMap
	_ = json.Unmarshal(resp, &image)
	//cache.RedisClient.SetMap(g.tag + consts.Images, image.Id, image)
	log.Printf("==============Set image %s visibility %s success", image.Id, property)
	return image.Id
}

func (g *Glance) CreateImageMember(imageId, memberId string) {
	createSuffix := fmt.Sprintf("/images/%s/members", imageId)
	formatter := `{"member": "%+v"}`
	reqBody := fmt.Sprintf(formatter, memberId)
	resp := g.Post(g.headers, createSuffix, reqBody)
	var imageMember entity.ImageMember
	_ = json.Unmarshal(resp, &imageMember)
	//cache.RedisClient.SetMap(g.tag + consts.Images, image.Id, image)
	log.Println("==============Create image member success", imageMember.ImageId)
}

func (g *Glance) SetImageMemberStatus(imageId, memberId, status string) {
	// status must in [pending, accepted, rejected]
	createSuffix := fmt.Sprintf("/images/%s/members/%s", imageId, memberId)
	formatter := `{"status": "%+v"}`
	reqBody := fmt.Sprintf(formatter, status)
	resp := g.Put(g.headers, createSuffix, reqBody)
	var imageMember entity.ImageMember
	_ = json.Unmarshal(resp, &imageMember)
	//cache.RedisClient.SetMap(g.tag + consts.Images, image.Id, image)
	log.Println("==============Set image status success", imageMember.ImageId)
}

func (g *Glance) DeleteImageMember(imageId, memberId string) {
	createSuffix := fmt.Sprintf("/images/%s/members/%s", imageId, memberId)
	if ok, _ := g.Delete(g.headers, createSuffix); ok {
		log.Println("==============Delete image member success", memberId)
	} else {
		log.Println("==============Delete image member failed", memberId)
	}
}

func (g *Glance) GetImage(imageId string) entity.ImageMap {
	suffix := fmt.Sprintf("/images/%s", imageId)
	resp := g.Get(g.headers, suffix)
	var image entity.ImageMap
	_ = json.Unmarshal(resp, &image)
	log.Println("==============Get image success")
	return image
}

func (g *Glance) GetImages() entity.Images {
	suffix := fmt.Sprintf("/images?project_id=%s", g.projectId)
	resp := g.Get(g.headers, suffix)
	var images entity.Images
	_ = json.Unmarshal(resp, &images)
	log.Println("==============List image success", images.Is)
	return images
}

func (g *Glance) ConstructRawImage() string {
	containerFormat := "bare"
	diskFormat := "raw"
    name := "dx_image"
    formatter := `{
                     "container_format": "%+v",
                     "disk_format": "%+v",
                     "name": "%+v",
                 }`
    body := fmt.Sprintf(formatter, containerFormat, diskFormat, name)
    return body
}

func (g *Glance) ConstructS3Image() string {
	containerFormat := "bare"
	diskFormat := "raw"
	name := "dx_image"
	formatter := `{
                     "container_format": "%+v",
                     "disk_format": "%+v",
                     "name": "%+v",
                     "properties": [
                         "s3_store_access_key": "56P9PIKO63GIO70L7B4P",
                         "s3_store_secret_key": "LCtUoNr2dvuyDscvcQspoF6dEOmeqWIM4yfLCMyw",
                         "s3_store_bucket": "image_bucket",
                         "s3_store_region_name": "us-east-1",
                         "s3_store_host": "http://10.50.114.157:6780",
                         "s3_store_create_bucket_on_put": "True",
                         "s3_store_bucket_url_format": "auto",
                         "s3_store_large_object_size": "100",
                         "s3_store_large_object_chunk_size": "10",
                         "s3_store_thread_pools": "10"
                     ]
                 }`
	body := fmt.Sprintf(formatter, containerFormat, diskFormat, name)
	return body
}

func (g *Glance) CreateRawImage() string {
    reqBody := g.ConstructRawImage()
    return g.CreateImage(reqBody)
}

func (g *Glance) CreateImageToS3() string {
	g.headers["OpenStack-image-store-ids"] = "cheap"
	reqBody := g.ConstructS3Image()
	return g.CreateImage(reqBody)
}

func (g *Glance) DeleteImage(imageId string) {
	urlSuffix := "/images/" + imageId
	if ok, _ := g.Delete(g.headers, urlSuffix); ok {
		log.Println("==============Delete image success", imageId)
	} else {
		log.Println("==============Delete image failed", imageId)
	}
}

func (g *Glance) DeleteImages() {
	images := g.GetImages()
	for _, image := range images.Is {
		g.DeleteImage(image.Id)
	}
}

func (g *Glance) GetImageSchemas() {
	suffix := "/schemas/images"
	resp := g.Get(g.headers, suffix)
	fmt.Println(resp)
	//images := resp["images"]
	//for _, image := range images.([]interface{}) {
	//	fmt.Println("image==", image.(map[string]interface{}))
	//}
}

// GetStores not implement
func (g *Glance) GetStores()  {
	suffix := "/info/stores"
	resp := g.Get(g.headers, suffix)
	fmt.Println(resp)
}

func (g *Glance) GetImportMethods()  {
	suffix := "/info/import"
	resp := g.DecorateGetResp(g.Get)(g.headers, suffix)
	importMethods := resp["import-methods"]
	fmt.Println(importMethods)
}

