package main

import (
	"github.com/therecipe/qt/core"
)

const (
	Name = int(core.Qt__UserRole) + 1<<iota
	Score
	DestValue
	Gender
)

type ResultModel struct {
	core.QAbstractListModel

	_ func() `constructor:"init"`

	_ map[int]*core.QByteArray `property:"roles"`
	_ []*Result                `property:"results"`

	_ func(string, string, float32, uint) `slot:"addResult"`
}

type Result struct {
	core.QObject

	_ string  `property:"name"`
	_ float32 `property:"score"`
	_ uint32  `property:"destValue"`
	_ string  `property:"gender"`
}

func init() {
	Result_QRegisterMetaType()
}

func (r *ResultModel) init() {
	r.SetRoles(map[int]*core.QByteArray{
		Name:      core.NewQByteArray2("name", len("name")),
		Score:     core.NewQByteArray2("score", len("score")),
		Gender:    core.NewQByteArray2("gender", len("gender")),
		DestValue: core.NewQByteArray2("destValue", len("destValue")),
	})

	r.ConnectRowCount(r.rowCount)
	r.ConnectColumnCount(r.columnCount)
	r.ConnectData(r.data)
	r.ConnectRoleNames(r.roleNames)
	r.ConnectAddResult(r.addResult)
	// r.ConnectModelReset(r.modelReset)
}

func (r *ResultModel) rowCount(parent *core.QModelIndex) int {
	return len(r.Results())
}

func (r *ResultModel) columnCount(parent *core.QModelIndex) int {
	return 1
}

func (r *ResultModel) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() || index.Row() > r.rowCount(index) {
		return core.NewQVariant()
	}

	result := r.Results()[index.Row()]

	switch role {
	case Name:
		return core.NewQVariant14(result.Name())
	case Gender:
		return core.NewQVariant14(result.Gender())
	case Score:
		return core.NewQVariant13(result.Score())
	case DestValue:
		return core.NewQVariant8(result.DestValue())
	default:
		return core.NewQVariant()
	}
}

func (r *ResultModel) roleNames() map[int]*core.QByteArray {
	return r.Roles()
}

func (r *ResultModel) modelReset() {
	r.BeginResetModel()
	r.SetResults(nil)
	r.EndResetModel()
}

func (r *ResultModel) addResult(playerName, gender string, score float32, destValue uint) {
	r.BeginInsertRows(core.NewQModelIndex(), len(r.Results()), len(r.Results()))

	result := NewResult(nil)
	result.SetScore(score)
	result.SetName(playerName)
	result.SetGender(gender)
	result.SetDestValue(destValue)
	r.SetResults(append(r.Results(), result))

	r.EndInsertRows()
}
