package software

import (
	"os/exec"
	"time"

	"github.com/bellh14/DesignManager/pkg/utils/log"
)

func InstallSoftware(node, srcPath, tarBall, dest string, logger *log.Logger) error {
	// copy tarball into node's /temp directory

	nodeDest := node + ":" + dest
	mkdirCmd := exec.Command("ssh", node, "mkdir", dest)
	if err := mkdirCmd.Run(); err != nil {
		logger.Fatal("Failed to make dest directory", err)
	}

	time.Sleep(2 * time.Second)

	fullTarPath := srcPath + "/" + tarBall
	fullDestPath := nodeDest + "/" + tarBall
	scpCmd := exec.Command("scp", fullTarPath, fullDestPath)
	if err := scpCmd.Run(); err != nil {
		logger.Fatal("Failed to scp tar ball", err)
	}

	tarPath := dest + "/" + tarBall
	tarCmd := exec.Command("ssh", node, "tar", "-xf", tarPath, "-C", dest)

	err := tarCmd.Run()
	if err != nil {
		logger.Fatal("Failed to install software", err)
	}

	// err := os.MkdirAll(dest, 0o777)
	// if err != nil {
	// 	logger.Fatal("Failed to make dest directory", err)
	// }
	//
	// err = utils.CopyFile(tarBall, dest)
	// if err != nil {
	// 	logger.Fatal(
	// 		fmt.Sprintf("Failed to copy software tar ball: %s, into %s", tarBall, dest),
	// 		err,
	// 	)
	// 	return err
	// }
	//
	// cmd := exec.Command("tar", "-xf", tarBall, "-C", dest)
	// cmd.Stderr = os.Stderr
	//
	// if err = cmd.Run(); err != nil {
	// 	logger.Fatal("Failed to extract tarball: ", err)
	// 	return err
	// }

	return nil
}
