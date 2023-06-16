package md

import (
	"encoding/json"
	"testing"

	"github.com/woweh/forge-api-go-client/md"
)

func Test_IfcAdvancedSpec_Json(t *testing.T) {
	type args struct {
		conversionMethod md.ConversionMethod
		storeys          md.Option
		spaces           md.Option
		openings         md.Option
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "All params are filled in",
			args: args{conversionMethod: md.V3, storeys: md.Hide, spaces: md.Show, openings: md.Skip},
			want: "{\"conversionMethod\":\"v3\",\"buildingStoreys\":\"hide\",\"spaces\":\"show\",\"openingElements\":\"skip\"}",
		},
		{
			name: "Legacy method - no additional parameters",
			args: args{conversionMethod: md.Legacy},
			want: "{\"conversionMethod\":\"legacy\"}",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				advancedSpec := md.IfcAdvancedSpec(
					tt.args.conversionMethod, tt.args.storeys, tt.args.spaces, tt.args.openings,
				)
				bytes, _ := json.Marshal(advancedSpec)
				gotJson := string(bytes)
				if gotJson != tt.want {
					t.Errorf("IfcAdvancedSpec() json = %v, want %v", gotJson, tt.want)
				}
			},
		)
	}
}

func Test_RevitAdvancedSpec_Json(t *testing.T) {
	type args struct {
		generateMasterViews bool
		materialMode        md.MaterialMode
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "All params are filled in",
			args: args{true, md.Auto},
			want: "{\"generateMasterViews\":true,\"materialMode\":\"auto\"}",
		},
		{
			name: "Only generateMasterViews, no materialMode",
			args: args{false, ""},
			want: "{\"generateMasterViews\":false}",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				advancedSpec := md.RevitAdvancedSpec(&tt.args.generateMasterViews, tt.args.materialMode)
				bytes, _ := json.Marshal(advancedSpec)
				gotJson := string(bytes)
				if gotJson != tt.want {
					t.Errorf("RevitAdvancedSpec() json = %v, want %v", gotJson, tt.want)
				}
			},
		)
	}
}

func TestNavisworksAdvancedSpec(t *testing.T) {
	type args struct {
		hiddenObjects              *bool
		basicMaterialProperties    *bool
		autodeskMaterialProperties *bool
		timeLinerProperties        *bool
	}
	trueValue := true
	falseValue := false
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "All params are filled in",
			args: args{&trueValue, &falseValue, &trueValue, &falseValue},
			want: "{\"hiddenObjects\":true,\"basicMaterialProperties\":false,\"autodeskMaterialProperties\":true,\"timelinerProperties\":false}",
		},
		{
			name: "Only hiddenObjects parameter",
			args: args{&trueValue, nil, nil, nil},
			want: "{\"hiddenObjects\":true}",
		},
		{
			name: "Only hiddenObjects and timelinerProperties parameters",
			args: args{&trueValue, nil, nil, &falseValue},
			want: "{\"hiddenObjects\":true,\"timelinerProperties\":false}",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				advancedSpec := md.NavisworksAdvancedSpec(
					tt.args.hiddenObjects, tt.args.basicMaterialProperties, tt.args.autodeskMaterialProperties,
					tt.args.timeLinerProperties,
				)
				bytes, _ := json.Marshal(advancedSpec)
				gotJson := string(bytes)
				if gotJson != tt.want {
					t.Errorf("NavisworksAdvancedSpec() json = %v, want %v", gotJson, tt.want)
				}
			},
		)
	}
}

func Test_ObjAdvancedSpec_Json(t *testing.T) {
	type args struct {
		exportFileStructure md.ExportFileStructure
		unit                md.Unit
		modelGuid           string
		objectIds           []int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "All params are filled in",
			args: args{
				exportFileStructure: md.Single, unit: md.Meter, modelGuid: "justSomeGuid", objectIds: []int{1, 2, 3},
			},
			want: "{\"exportFileStructure\":\"single\",\"unit\":\"meter\",\"modelGuid\":\"justSomeGuid\",\"objectIds\":[1,2,3]}",
		},
		{
			name: "Minimum params are filled in",
			args: args{exportFileStructure: "", unit: md.None, modelGuid: "justSomeGuid", objectIds: []int{1, 2, 3}},
			want: "{\"modelGuid\":\"justSomeGuid\",\"objectIds\":[1,2,3]}",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				advancedSpec := md.ObjAdvancedSpec(
					tt.args.exportFileStructure, tt.args.unit, tt.args.modelGuid, &tt.args.objectIds,
				)
				bytes, _ := json.Marshal(advancedSpec)
				gotJson := string(bytes)
				if gotJson != tt.want {
					t.Errorf("ObjAdvancedSpec() json = %v, want %v", gotJson, tt.want)
				}
			},
		)
	}
}

func Test_AdvancedSpec_IsEmpty(t *testing.T) {

	// if no AdvancedSpec is defined, it shouldn't be in the json
	expectedJson := "{\"type\":\"svf2\",\"views\":[\"3d\"]}"

	formatSpec := md.FormatSpec{
		Type:  md.SVF2,
		Views: []md.ViewType{md.View3D},
	}
	bytes, _ := json.Marshal(formatSpec)
	gotJson := string(bytes)

	if gotJson != expectedJson {
		t.Errorf("formatSpec() json = %v, want %v", gotJson, expectedJson)
	}
}
