package ai

import (
	"testing"
)

func mockFamilys() []Family {
	albertAis := []*Ai{
		{
			Name:        "Albert",
			Version:     0,
			Description: "",
			FileName:    "AIAlbert.cc",
		},
		{
			Name:        "Albert",
			Version:     1,
			Description: "",
			FileName:    "AIAlbert_1.cc",
		},
		{
			Name:        "Albert",
			Version:     2,
			Description: "avoid_enemy",
			FileName:    "AIAlbert_2_avoid_enemy.cc",
		},
		{
			Name:        "Albert",
			Version:     3,
			Description: "",
			FileName:    "AIAlbert_3.cc",
		},
	}

	demo := Ai{
		Name:        "Demo",
		Version:     0,
		Description: "",
		FileName:    "AIDemo.cc",
	}

	dummy := Ai{
		Name:        "Dummy",
		Version:     0,
		Description: "",
		FileName:    "AIDummy.cc",
	}

	albertFamily := Family{}
	demoFamily := Family{}
	dummyFamily := Family{}

	albertFamily.Add(albertAis...)
	demoFamily.Add(&demo)
	dummyFamily.Add(&dummy)

	return []Family{
		albertFamily,
		demoFamily,
		dummyFamily,
	}
}

func TestMatchesAnyDescriptors(t *testing.T) {
	familys := mockFamilys()

	tests := []struct {
		Ai          *Ai
		Descriptors []string
		Matches     bool
	}{
		{familys[0].Ais[0], []string{"Albert"}, false},
		{familys[0].Ais[0], []string{"Albert:0"}, true},
		{familys[0].Ais[0], []string{"Albert:-4"}, true},
		{familys[0].Ais[0], []string{"Albert:"}, true},
		{familys[0].Ais[0], []string{"Albert:0..3"}, true},
		{familys[0].Ais[0], []string{"Albert:1..3"}, false},
		{familys[0].Ais[0], []string{"Albert:0..1"}, true},
		{familys[0].Ais[0], []string{"Albert:0.."}, true},
		{familys[0].Ais[0], []string{"Albert:.."}, true},

		{familys[0].Ais[3], []string{"Albert"}, true},
		{familys[0].Ais[3], []string{"Albert:-1"}, true},
		{familys[0].Ais[2], []string{"Albert"}, false},
		{familys[0].Ais[2], []string{"Albert:-2"}, true},
		{familys[0].Ais[2], []string{"Albert:..-2"}, true},

		{familys[1].Ais[0], []string{"Demo:"}, true},
		{familys[1].Ais[0], []string{"Demo"}, true},
		{familys[1].Ais[0], []string{"Albert:"}, false},
	}

	for _, test := range tests {
		descriptors := make([]Descriptor, 0)

		for _, descriptorString := range test.Descriptors {
			descriptors = append(descriptors, DescriptorFromString(descriptorString))
		}

		if matchesAnyDescriptors(test.Ai, descriptors...) != test.Matches {
			t.Errorf("Failed for %s, %s, %t", test.Ai.FileName, test.Descriptors[0], test.Matches)
		}
	}
}

func TestGetAis(t *testing.T) {
	tests := []struct {
		FileNames []string
		Ais       []struct {
			Ai          Ai
			LastVersion int
		}
	}{
		{
			FileNames: []string{
				"api.pdf",
				"AICancellara.cc",
				"board.cc",
				"AIDemo_1.cc",
				"AICancellara_1.cc",
				"AIAlbert_3_description.cc",
			},
			Ais: []struct {
				Ai          Ai
				LastVersion int
			}{
				{
					Ai: Ai{
						Name:        "Cancellara",
						Version:     0,
						Description: "",
						FileName:    "AICancellara.cc",
					},
					LastVersion: 1,
				},
				{
					Ai: Ai{
						Name:        "Demo",
						Version:     1,
						Description: "",
						FileName:    "AIDemo_1.cc",
					},
					LastVersion: 1,
				},
				{
					Ai: Ai{
						Name:        "Cancellara",
						Version:     1,
						Description: "",
						FileName:    "AICancellara_1.cc",
					},
					LastVersion: 1,
				},
				{
					Ai: Ai{
						Name:        "Albert",
						Version:     3,
						Description: "description",
						FileName:    "AIAlbert_3_description.cc",
					},
					LastVersion: 3,
				},
			},
		},
	}

	for _, test := range tests {
		ais := GetAis(test.FileNames)

		for i, ai := range ais {
			if ai.Name != test.Ais[i].Ai.Name {
				t.Error("Found Name", ai.Name, ", expected", test.Ais[i].Ai.Name)
			}

			if ai.Version != test.Ais[i].Ai.Version {
				t.Error("Found Version", ai.Version, ", expected", test.Ais[i].Ai.Version)
			}

			if ai.Description != test.Ais[i].Ai.Description {
				t.Error("Found Description", ai.Description, ", expected", test.Ais[i].Ai.Description)
			}

			if ai.FileName != test.Ais[i].Ai.FileName {
				t.Error("Found FileName", ai.FileName, ", expected", test.Ais[i].Ai.FileName)
			}

			if ai.Family.Name != test.Ais[i].Ai.Name {
				t.Error("Found FamilyName", ai.Family.Name, ", expected", test.Ais[i].Ai.Name)
			}

			if ai.Family.LastVersion.Version != test.Ais[i].LastVersion {
				t.Error("For ", ai, "Found LastVersion", ai.Family.LastVersion.Version, ", expected", test.Ais[i].LastVersion)
			}
		}

	}

}
