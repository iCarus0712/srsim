package onlysilenceremains

import (
	"github.com/simimpact/srsim/pkg/engine"
	"github.com/simimpact/srsim/pkg/engine/equip/lightcone"
	"github.com/simimpact/srsim/pkg/engine/event"
	"github.com/simimpact/srsim/pkg/engine/info"
	"github.com/simimpact/srsim/pkg/engine/modifier"
	"github.com/simimpact/srsim/pkg/engine/prop"
	"github.com/simimpact/srsim/pkg/key"
	"github.com/simimpact/srsim/pkg/model"
)

const (
	OnlySilenceRemains       key.Modifier = "only_silence_remains"
	OnlySilenceRemainsCRBuff key.Modifier = "only_silence_remians_cr_buff"
)

// Increases ATK of its wearer by 16/20/24/28/32%. If there are 2 or fewer
// enemies on the field, increases wearer's CRIT Rate by 12/15/18/21/24%.
func init() {
	lightcone.Register(key.OnlySilenceRemains, lightcone.Config{
		CreatePassive: Create,
		Rarity:        4,
		Path:          model.Path_HUNT,
		Promotions:    promotions,
	})

	modifier.Register(OnlySilenceRemains, modifier.Config{})
	modifier.Register(OnlySilenceRemainsCRBuff, modifier.Config{})
}

// Note: does not properly handle enemies running away (such as trotters)
func Create(engine engine.Engine, owner key.TargetID, lc info.LightCone) {
	atk_amt := 0.12 + 0.04*float64(lc.Imposition)
	cr_amt := 0.09 + 0.03*float64(lc.Imposition)

	engine.Events().EnemyAdded.Subscribe(func(e event.EnemyAddedEvent) {
		updateCRBuff(engine, owner, cr_amt)
	})

	engine.Events().TargetDeath.Subscribe(func(e event.TargetDeathEvent) {
		updateCRBuff(engine, owner, cr_amt)
	})

	engine.AddModifier(owner, info.Modifier{
		Name:   OnlySilenceRemains,
		Source: owner,
		Stats:  info.PropMap{prop.ATKPercent: atk_amt},
	})
}

func updateCRBuff(engine engine.Engine, owner key.TargetID, amt float64) {
	if len(engine.Enemies()) <= 2 {
		engine.AddModifier(owner, info.Modifier{
			Name:   OnlySilenceRemainsCRBuff,
			Source: owner,
			Stats:  info.PropMap{prop.CritChance: amt},
		})
	} else {
		engine.RemoveModifier(owner, OnlySilenceRemainsCRBuff)
	}
}