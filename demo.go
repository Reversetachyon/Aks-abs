package main

import (
	"fmt"     //for showing output
	"net/http"
	"os" // for opening file

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob" // azureblob sdk
	"github.com/gin-gonic/gin"                             // gin framework
)

func main() {
	// get connectionstring from azure control panel in access key of storage account
	connectionString := "DefaultEndpointsProtocol=https;AccountName=testcontainerintern;AccountKey=PumBxb7dhYRqjoOfaG9jeQ5Mfe28thj1pXkC0NmWOANWXv9GjKDQCSaG5o/XXhqvAwoi17o9ZP0984I26c5y5g==;EndpointSuffix=core.windows.net"

	//create container client for manipulating container and blobs
	containerClient,_:=azblob.NewContainerClientFromConnectionString(connectionString,"hellotest",&azblob.ClientOptions{})

	//create server
	server:=gin.New()

	//Upload file
	server.GET("/upload",func(ctx *gin.Context) {
		UploadToBlobStotage(containerClient,ctx)
	})
	
	//list blob in container
	server.GET("/list",func(ctx *gin.Context) {
		ListBlobInContainer(containerClient,ctx)
	})

	//Download file
	server.GET("/download",func(ctx *gin.Context) {
		DownloadBlob(containerClient,ctx,"test.json")
	})
	
	//Delete file
	server.GET("/delete",func(ctx *gin.Context) {
		DeleteBlob(containerClient,ctx,"test.json")
	})
	
	server.Run(":8080")

}

//Upload file to container in azure blob storage
func UploadToBlobStotage(contClient azblob.ContainerClient, ctx *gin.Context ) {
	// Open file 
	file, err := os.Open("test.json")
	// handle error when opening file
	if err != nil {
		fmt.Println(err.Error())
	}
	bc:=contClient.NewBlockBlobClient(file.Name())
	resp,err := bc.UploadFileToBlockBlob(ctx,file,azblob.HighLevelUploadToBlockBlobOption{})
	// handle error when upload file to container
	if err != nil {
		fmt.Println(err.Error())
	}else{
		// show statuscode when upload file complete
		fmt.Println(resp.StatusCode)
		ctx.JSON(http.StatusAccepted,map[string]interface{}{
			"Response":resp.StatusCode,
		})
	}

}

//List files in container
func ListBlobInContainer(contClient azblob.ContainerClient,  ctx *gin.Context){
	lbf := contClient.ListBlobsFlat(&azblob.ContainerListBlobFlatSegmentOptions{})
	resp := make([]string,0,5)
	for lbf.NextPage(ctx){
		bresp := lbf.PageResponse()
		for _,response := range bresp.ListBlobsFlatSegmentResponse.Segment.BlobItems{
			fmt.Println(*response.Name) // show every file name
			resp = append(resp, *response.Name)
		}
	}
	ctx.JSON(http.StatusAccepted,gin.H{
		"Response": resp,
	})
}

//Download file in container
func DownloadBlob(contClient azblob.ContainerClient,  ctx *gin.Context, fileName string){
	// Reserve buffer to contain downloaded data
	var blob = make([]byte,100)
	bc:=contClient.NewBlockBlobClient(fileName)
   	err :=bc.DownloadBlobToBuffer(ctx,0,0,blob,azblob.HighLevelDownloadFromBlobOptions{})
   		//handle error of download
   		if err != nil {
	   		fmt.Println(err.Error())
   		}
	// show data
   	fmt.Println(string(blob))
	ctx.JSON(http.StatusAccepted,gin.H{
		"Response": string(blob),
	})
}

//Delete file in container
func DeleteBlob(contClient azblob.ContainerClient,  ctx *gin.Context, fileName string){
	bc:=contClient.NewBlockBlobClient(fileName)
	resp,err := bc.Delete(ctx,&azblob.DeleteBlobOptions{})
		//handle error of delete
		if err != nil {
			fmt.Println(err.Error())
		}
	// show status when delete completed
	fmt.Println(resp.RawResponse.StatusCode)
	ctx.JSON(http.StatusAccepted,gin.H{
		"Response": resp.RawResponse.StatusCode,
	})
}