package models

type (
	DirectoryEntry struct {
		Permissions string
		Owner       string
		Group       string
		Size        int64
		Modified    string
		Name        string
	}

	DownloadFile struct {
		TargetPath string `form:"targetpath" json:"targetpath" binding:"required"`
		SavePath   string `form:"savepath" json:"savepath" binding:"required"`
	}

	UploadFile struct {
		TargetPath    string `form:"targetpath" json:"targetpath"`
		LocalFilePath string `form:"localfilepath" json:"localfilepath" binding:"required"`
	}
)
