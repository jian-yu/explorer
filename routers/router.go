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

	"github.com/astaxie/beego"
)

func init() {

	mgoStore := db.NewMongoStore()
	baseController := controllers.NewBaseController(mgoStore)

	krc := &account.KindsRewardController{Base: baseController}

	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/public",
			beego.NSInclude(
				&controllers.PublicController{
					Base: baseController,
				},
			),
		),
		beego.NSNamespace("/drawing",
			beego.NSInclude(
				&controllers.DrawingDataController{Base: baseController},
			),
		),
		beego.NSNamespace("/blocks",
			beego.NSInclude(
				&controllers.BlockController{Base: baseController},
			),
		),
		beego.NSNamespace("/txs",
			beego.NSInclude(
				&controllers.TxsController{Base: baseController},
			),
		),
		beego.NSNamespace("/tx",
			beego.NSInclude(
				&controllers.TxDetailControllers{Base: baseController},
			),
		),
		beego.NSNamespace("/blockTxs",
			beego.NSInclude(
				&controllers.BlockTxController{Base: baseController},
			),
		),
		beego.NSNamespace("/validators",
			beego.NSInclude(
				&controllers.ValidatorsController{Base: baseController},
			),
			//beego.NSRouter("/", &controllers.ValidatorsController{}),
		),
		beego.NSNamespace("/validatorBase",
			beego.NSInclude(
				&validator.VaBaseInfoController{Base: baseController},
			),
		),
		beego.NSNamespace("/validatorDelegations",
			beego.NSInclude(
				&validator.DelegationsController{Base: baseController},
			),
		),
		beego.NSNamespace("/validatorPowerEvent",
			beego.NSInclude(
				&validator.PowerEventController{Base: baseController},
			),
		),
		beego.NSNamespace("/validatorProposedBlock",
			beego.NSInclude(
				&validator.ProposedBlocksController{Base: baseController},
			),
		),
		beego.NSNamespace("/accountInfo",
			beego.NSInclude(
				&account.BaseInfoController{Base: baseController, KindsRewardController: krc},
			),
		),
		beego.NSNamespace("/delegators",
			beego.NSInclude(
				&account.DeleatorsController{Base: baseController},
			),
		),
		beego.NSNamespace("/unbonding",
			beego.NSInclude(
				&account.UnbondingsController{Base: baseController},
			),
		),
		beego.NSNamespace("/delegatorTxs",
			beego.NSInclude(
				&account.DelegatorTxController{Base: baseController},
			),
		),
		beego.NSNamespace("/delegatorAllKindsReward",
			beego.NSInclude(krc),
		),
	)

	beego.AddNamespace(ns)
}
