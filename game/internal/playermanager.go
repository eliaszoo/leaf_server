package internal

type playerManger struct {
	PlayerMap map[string]*Player
}

func NewPlayerManager() *playerManger {
	mgr := new(playerManger)
	mgr.Init()
	return mgr
}

func (self *playerManger) Init() {
	self.PlayerMap = make(map[string]*Player)
}

func (self *playerManger) Get(id string) *Player {
	return self.PlayerMap[id]
}

func (self *playerManger) AddPlayer(player *Player) {
	self.PlayerMap[player.objid] = player
}

func (self *playerManger) DelPlayer(objid string) {
	delete(self.PlayerMap, objid)
}

func (self *playerManger) Close() {
	for _, player := range self.PlayerMap {
		player.SaveSync()
	}
}