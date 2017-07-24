package gologo

type ScreenDirection int

const (
    SCREEN_UP ScreenDirection = iota
    SCREEN_DOWN
    SCREEN_LEFT
    SCREEN_RIGHT
)

func PhysicsTick() {
   ResolveContacts(GenerateContacts())
}
