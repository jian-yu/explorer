// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"explorer/controllers"
	"explorer/controllers/account"
	"explorer/controllers/validator"
	"explorer/db"
	"explorer/utils"
	"github.com/astaxie/beego"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func init() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	configPath := os.Getenv(`explorer_config_path`)

	_, err := utils.LoadViper("", configPath)
	if err != nil {
		panic(err)
	}

	mgoStore := db.NewMongoStore()
	baseController := controllers.NewBaseController(mgoStore)

	krc := &account.KindsRewardController{Base: baseController}

	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&controllers.PublicController{
					Base: baseController,
				},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&controllers.DrawingDataController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&controllers.BlockController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&controllers.TxsController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&controllers.TxDetailControllers{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&controllers.BlockTxController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&controllers.ValidatorsController{Base: baseController},
			),
			//beego.NSRouter("/", &controllers.ValidatorsController{}),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&validator.VaBaseInfoController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&validator.DelegationsController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&validator.PowerEventController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&validator.ProposedBlocksController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&account.BaseInfoController{Base: baseController, KindsRewardController: krc},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&account.DeleatorsController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&account.UnbondingsController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(
				&account.DelegatorTxController{Base: baseController},
			),
		),
		beego.NSNamespace("/v1",
			beego.NSInclude(krc),
		),
	)

	beego.AddNamespace(ns)
}
