package md

import (
	"encoding/json"
	"testing"

	"github.com/woweh/forge-api-go-client/md"
)

func Test_IfcAdvancedSpec_Json(t *testing.T) {
	type args struct {
		conversionMethod md.IfcConversionMethod
		storeys          md.IfcOption
		spaces           md.IfcOption
		openings         md.IfcOption
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "All params are filled in",
			args: args{conversionMethod: md.IfcV3, storeys: md.IfcHide, spaces: md.IfcShow, openings: md.IfcSkip},
			want: "{\"conversionMethod\":\"v3\",\"buildingStoreys\":\"hide\",\"spaces\":\"show\",\"openingElements\":\"skip\"}",
		},
		{
			name: "IfcLegacy method - no additional parameters",
			args: args{conversionMethod: md.IfcLegacy},
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
		materialMode        md.RvtMaterialMode
		twoDViews           md.Rvt2dViews
		version             md.RvtExtractorVersion
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "All params are filled in",
			args: args{true, md.RvtAuto, md.RvtPdf, md.RvtNext},
			want: "{\"2dviews\":\"pdf\",\"extractorVersion\":\"next\",\"generateMasterViews\":true,\"materialMode\":\"auto\"}",
		},
		{
			name: "Only generateMasterViews",
			args: args{false, "", "", ""},
			want: "{\"generateMasterViews\":false}",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				advancedSpec := md.RevitAdvancedSpec(
					tt.args.generateMasterViews, tt.args.materialMode, tt.args.twoDViews, tt.args.version,
				)
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
		hiddenObjects              bool
		basicMaterialProperties    bool
		autodeskMaterialProperties bool
		timeLinerProperties        bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "True False True False",
			args: args{true, false, true, false},
			want: "{\"hiddenObjects\":true,\"basicMaterialProperties\":false,\"autodeskMaterialProperties\":true,\"timelinerProperties\":false}",
		},
		{
			name: "All True",
			args: args{true, true, true, true},
			want: "{\"hiddenObjects\":true,\"basicMaterialProperties\":true,\"autodeskMaterialProperties\":true,\"timelinerProperties\":true}",
		},
		{
			name: "All False",
			args: args{false, false, false, false},
			want: "{\"hiddenObjects\":false,\"basicMaterialProperties\":false,\"autodeskMaterialProperties\":false,\"timelinerProperties\":false}",
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
		exportFileStructure md.ObjExportFileStructure
		unit                md.ObjUnit
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
				exportFileStructure: md.ObjSingle, unit: md.ObjMeter, modelGuid: "justSomeGuid",
				objectIds: []int{1, 2, 3},
			},
			want: "{\"exportFileStructure\":\"single\",\"unit\":\"meter\",\"modelGuid\":\"justSomeGuid\",\"objectIds\":[1,2,3]}",
		},
		{
			name: "Minimum params are filled in",
			args: args{exportFileStructure: "", unit: md.ObjNone, modelGuid: "justSomeGuid", objectIds: []int{1, 2, 3}},
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
	expectedJson := "{\"type\":\"svf\",\"views\":[\"3d\"]}"

	formatSpec := md.FormatSpec{
		Type:  md.SVF,
		Views: md.ViewType3D(),
	}
	bytes, _ := json.Marshal(formatSpec)
	gotJson := string(bytes)

	if gotJson != expectedJson {
		t.Errorf("formatSpec() json = %v, want %v", gotJson, expectedJson)
	}
}
