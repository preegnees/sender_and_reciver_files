package distribution

import (
	"io"
	"log"
	"fmt"

	common "files/pkg/common"
	controller "files/pkg/reciver/writer"
)

type FileDistributor struct{
	cntrl *controller.ControllerOfFiles
}

func New(controller *controller.ControllerOfFiles) *FileDistributor {
	return &FileDistributor{
		cntrl: controller,
	}
}

var _ io.Writer = (*FileDistributor)(nil)

func (r *FileDistributor) Write(buf []byte) (n int, err error) {
	fmt.Printf("%s\n", buf)
	offset, path, index, sizeOpt, sizeFile, sizeBuffer := common.DeconstructOptions(buf)
	log.Println("offset:", offset)
	log.Println("path:", path)
	log.Println("index:", index)
	log.Println("sizeOpt:", sizeOpt)
	log.Println("sizeFile:", sizeFile)
	log.Println("sizeBuffer:", sizeBuffer)
	log.Printf("message: %s\n", buf[:len(buf)-int(sizeOpt)])
	n, err = r.cntrl.WriteToFile(r.cntrl.ParentDir + path, buf[:len(buf)-int(sizeOpt)])
	log.Println(err)
	if err != nil {
		return n, err
	}
	return len(buf), nil
}