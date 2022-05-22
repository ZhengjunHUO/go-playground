package main

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
        ut "github.com/go-playground/universal-translator"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
)

type k8sCluster struct {
	Name		string	`validate:"required,validateClusterName"`
	CNI		string	`validate:"required,oneof=cilium calico flannel weave"`
	Size		int	`validate:"required,numeric,min=1"`
	IsManaged	bool	`validate:"omitempty"`
	IsBaremetal	bool	`validate:"omitempty"`
	IsOverlay	bool	`validate:"omitempty"`
}

// 自定义tag对应的处理逻辑
func validateClusterName(fl validator.FieldLevel) bool {
	clusterName := fl.Field().String()	

	// 字段长度需介于3-20个字符之间
	if len(clusterName) < 3 || len(clusterName) > 20 {
		return false
	}
	
	// 合法字段名格式为字母下划线字母xx_xxx
	rgx, _ := regexp.Compile("[[:alpha:]]+_[[:alpha:]]+")

	return rgx.MatchString(clusterName)
}

func main() {
	// 初始化validator
	vldt := validator.New()
	// 关联自定义tag和其handler func
	vldt.RegisterValidation("validateClusterName", validateClusterName)

	// 初始化翻译器
	english := en.New()
	zhongwen := zh.New()
	uni := ut.New(english, zhongwen, english)
	trans, _ := uni.GetTranslator("zh")

	// 绑定翻译器到validator
	_ = zhTrans.RegisterDefaultTranslations(vldt, trans)

	clusters := []k8sCluster{
		k8sCluster{
			Name:		"x",
			Size:		0,
			CNI:		"awsvps",
			IsManaged:	false,
		},
		k8sCluster{
			Name:		"huo_k8s",
			Size:		8,
			CNI:		"cilium",
			IsManaged:	false,
			IsBaremetal:	true,
			IsOverlay:	false,
		},
	}

	for i := range clusters {
		fmt.Printf("<Cluster %d>\n", i)
		// 检查实例是否满足struct中的定义
		if err := vldt.Struct(clusters[i]); err != nil {
			// 在valid返回的错误集中使用翻译器解读错误信息
			translatedMap := err.(validator.ValidationErrors).Translate(trans)
			// 打印错误信息
			for k, v := range translatedMap {
				fmt.Printf("%s: %s\n", k, v)
			}
		}else{
			fmt.Println("Validation ok.")
		}
	}
}
