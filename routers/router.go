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
	basecontroller := controllers.NewBaseController(mgoStore)

	krc := &account.KindsRewardController{Base: basecontroller}

	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/public",
			beego.NSInclude(
				&controllers.PublicController{
					Base: basecontroller,
				},
			),
		),
		beego.NSNamespace("/drawing",
			beego.NSInclude(
				&controllers.DrawingDataController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/blocks",
			beego.NSInclude(
				&controllers.BlockController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/txs",
			beego.NSInclude(
				&controllers.TxsController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/tx",
			beego.NSInclude(
				&controllers.TxDetailControllers{Base: basecontroller},
			),
		),
		beego.NSNamespace("/blockTxs",
			beego.NSInclude(
				&controllers.BlockTxController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/validators",
			beego.NSInclude(
				&controllers.ValidatorsController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/validatorBase",
			beego.NSInclude(
				&validator.VaBaseInfoController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/validatorDelegations",
			beego.NSInclude(
				&validator.DelegationsController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/validatorPowerEvent",
			beego.NSInclude(
				&validator.PowerEventController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/validatorProposedBlock",
			beego.NSInclude(
				&validator.ProposedBlocksController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/accountInfo",
			beego.NSInclude(
				&account.BaseInfoController{Base: basecontroller, KindsRewardController: krc},
			),
		),
		beego.NSNamespace("/delegators",
			beego.NSInclude(
				&account.DeleatorsController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/unbonding",
			beego.NSInclude(
				&account.UnbondingsController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/delegatorTxs",
			beego.NSInclude(
				&account.DelegatorTxController{Base: basecontroller},
			),
		),
		beego.NSNamespace("/delegatorAllKindsReward",
			beego.NSInclude(krc),
		),
	)

	beego.AddNamespace(ns)
}
