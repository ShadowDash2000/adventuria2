package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

const (
	ItemUseTypeInstant  = "instant"
	ItemUseTypeOnDrop   = "useOnDrop"
	ItemUseTypeOnReroll = "useOnReroll"
	ItemUseTypeOnRoll   = "useOnRoll"
)

type Effects struct {
	PointsIncrement int    `json:"pointsIncrement" mapstructure:"pointsIncrement"`
	JailEscape      bool   `json:"jailEscape" mapstructure:"jailEscape"`
	DiceMultiplier  int    `json:"diceMultiplier" mapstructure:"diceMultiplier"`
	DiceIncrement   int    `json:"diceIncrement" mapstructure:"diceIncrement"`
	Dices           []Dice `json:"dices" mapstructure:"changeDices"`
	IsSafeDrop      bool   `json:"isSafeDrop" mapstructure:"isSafeDrop"`
}

type Item struct {
	app     core.App
	item    *core.Record
	effects []*Effect
}

func NewItem(record *core.Record, app core.App) (*Item, error) {
	errs := app.ExpandRecord(record, []string{"effects"}, nil)
	if errs != nil {
		for _, err := range errs {
			return nil, err
		}
	}

	var effects []*Effect
	for _, effectRecord := range record.ExpandedAll("effects") {
		effect, err := NewEffect(effectRecord)
		if err != nil {
			return nil, err
		}

		effects = append(effects, effect)
	}

	item := &Item{
		app:     app,
		item:    record,
		effects: effects,
	}

	item.bindHooks()

	return item, nil
}

func (i *Item) bindHooks() {
	i.app.OnRecordAfterUpdateSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.item.Id {
			i.item = e.Record
		}
		return e.Next()
	})
}

func (i *Item) IsUsingSlot() bool {
	return i.item.GetBool("isUsingSlot")
}

func (i *Item) IsActiveByDefault() bool {
	return i.item.GetBool("isActiveByDefault")
}

func (i *Item) CanDrop() bool {
	return i.item.GetBool("canDrop")
}

func (i *Item) GetOrder() int {
	return i.item.GetInt("order")
}

func (i *Item) GetEffects() []*Effect {
	return i.effects
}

func (i *Item) GetEvent() string {
	return i.item.GetString("event")
}

func (i *Item) GetName() string {
	return i.item.GetString("name")
}
