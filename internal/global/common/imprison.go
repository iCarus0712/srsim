package common

import (
	"github.com/simimpact/srsim/pkg/engine/info"
	"github.com/simimpact/srsim/pkg/engine/modifier"
	"github.com/simimpact/srsim/pkg/engine/prop"
	"github.com/simimpact/srsim/pkg/key"
	"github.com/simimpact/srsim/pkg/model"
)

const (
	Imprisonment                   = "common-imprisonment"
	BreakImprisonment key.Modifier = "break-imprisonment"
)

type ImprisonState struct {
	SpeedDownRatio float64
	DelayRatio     float64
}

func init() {
	modifier.Register(Imprisonment, modifier.Config{
		Stacking:   modifier.ReplaceBySource,
		TickMoment: modifier.ModifierPhase1End,
		StatusType: model.StatusType_STATUS_DEBUFF,
		BehaviorFlags: []model.BehaviorFlag{
			model.BehaviorFlag_DISABLE_ACTION,
			model.BehaviorFlag_STAT_CTRL,
			model.BehaviorFlag_STAT_CONFINE,
			model.BehaviorFlag_STAT_SPEED_DOWN,
		},
		Listeners: modifier.Listeners{
			OnAdd: imprisonAdd,
		},
	})
}

func imprisonAdd(mod *modifier.Instance) {
	state, ok := mod.State().(ImprisonState)
	if !ok {
		panic("incorrect state used for imprisonment modifier")
	}

	mod.AddProperty(prop.SPDPercent, -state.SpeedDownRatio)
	mod.Engine().ModifyGaugeNormalized(info.ModifyAttribute{
		Key:    Imprisonment,
		Target: mod.Owner(),
		Source: mod.Source(),
		Amount: state.DelayRatio,
	})
}
