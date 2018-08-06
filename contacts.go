package gologo

import (
	"github.com/go-gl/mathgl/mgl32"
)

type PenetrationResolver interface {
	ResolvePenetration(contact *Contact)
}

type PostContactResolver interface {
	ResolveContact(contact *Contact)
}

type Contact struct {
	Objects               [2]*Object
	ContactGenerator      ContactGenerator
	GeneratorData         interface{}
	ContactNormal         mgl32.Vec4
	ContactPoint          mgl32.Vec4
	PenetrationResolver   PenetrationResolver
	Penetration           float32
	PenetrationResolution *PenetrationResolution
	PostContactResolver   PostContactResolver
}

type PenetrationResolution [2]mgl32.Vec3

type PerpendicularPenetrationResolver struct{}

type CallbackPostContactResolver struct {
	Callback func(*Contact)
}

func GenerateContacts() []*Contact {
	result := []*Contact{}

	for _, generator := range contactGenerators {
		contacts := generator.GenerateContacts()
		result = append(result, contacts...)
	}

	return result
}

func ResolveContacts(contacts []*Contact) {
	// Do one pass guaranteeing each contact is handles once
	for i := 0; i < len(contacts); i++ {
		if contacts[i].PenetrationResolver != nil {
			contacts[i].PenetrationResolver.ResolvePenetration(contacts[i])
			UpdateContacts(contacts, contacts[i])
		}
	}

	// Do a limited number of iterations resolving the deepest penetration
	// first.
	for i := 0; i < len(contacts)*MAX_CONTACT_ITERATIONS; i++ {
		var maxPenetration float32
		maxIdx := len(contacts)

		for j := 0; j < len(contacts); j++ {
			if contacts[j].Penetration > maxPenetration &&
				contacts[j].PenetrationResolver != nil &&
				contacts[j].PenetrationResolution != nil {
				maxIdx = j
			}
		}

		if maxIdx == len(contacts) {
			break
		}

		contacts[maxIdx].PenetrationResolver.ResolvePenetration(contacts[maxIdx])
		UpdateContacts(contacts, contacts[maxIdx])
	}

	// Invoke the post contact resolver for each contact
	for i := 0; i < len(contacts); i++ {
		if contacts[i].PostContactResolver != nil {
			contacts[i].PostContactResolver.ResolveContact(contacts[i])
		}
	}
}

func UpdateContacts(contacts []*Contact, resolved *Contact) {
	for i := 0; i < len(contacts); i++ {
		if contacts[i] != resolved {
			// Find contacts that share an object with the resolved contact
			for j := 0; j < len(resolved.Objects); j++ {
				for k := 0; k < len(contacts[i].Objects); k++ {
					if contacts[i].Objects[k] == resolved.Objects[j] {
						contacts[i].ContactGenerator.
							UpdateContact(contacts[i], k, resolved, j)
					}
				}
			}
		}
	}
}

func (p *PerpendicularPenetrationResolver) ResolvePenetration(contact *Contact) {
	contact.PenetrationResolution = &PenetrationResolution{}

	if contact.Penetration <= 0 {
		return
	}

	inverseMass0 := contact.Objects[0].Primitive.GetInverseMass()
	inverseMass1 := float32(0.0)
	totalInverseMass := inverseMass0

	if contact.Objects[1] != nil {
		inverseMass1 = contact.Objects[1].Primitive.GetInverseMass()
		totalInverseMass += inverseMass1
	}

	if totalInverseMass <= 0 {
		return
	}

	movePerInverseMass := contact.ContactNormal.Vec3()
	movePerInverseMass = movePerInverseMass.Mul(contact.Penetration / totalInverseMass)

	movement := movePerInverseMass
	movement = movement.Mul(inverseMass0)
	translation := contact.Objects[0].Model.Col(3).Add(movement.Vec4(0.0))
	contact.Objects[0].Model.SetCol(3, translation)
	contact.PenetrationResolution[0] = translation.Vec3()

	if contact.Objects[1] != nil {
		movement = movePerInverseMass
		movement = movement.Mul(-inverseMass1)
		translation = contact.Objects[1].Model.Col(3).Add(movement.Vec4(0.0))
		contact.Objects[1].Model.SetCol(3, translation)
		contact.PenetrationResolution[1] = translation.Vec3()
	}

	contact.ContactGenerator.UpdateContact(contact, 0, contact, 0)
}
