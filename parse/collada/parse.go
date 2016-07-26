package collada

import (
	"fmt"
	"encoding/xml"
	"os"
	"io/ioutil"
)
func Parse (fileName string) (error){
	relativePath, err := os.Getwd()
	if err != nil{
		return fmt.Errorf("parse collada: realative path read error: %v", err)
	}

	//Read entire collada file
	colladaBytes, err := ioutil.ReadFile(relativePath + "/../../../data/models/" + fileName)
	if err != nil{
		return fmt.Errorf("parse collada: %v", err)
	}
	var collada collada
	if err := xml.Unmarshal(colladaBytes, &collada); err != nil{
		return fmt.Errorf("parse collada: %v", err)
	}
	for _, sources := range collada.LibraryGeometries.Geometries[0].Meshes[0].Sources{
		fmt.Println(sources)
	}
	fmt.Printf("%#v\n", collada.LibraryGeometries.Geometries[0].Meshes[0])
	return nil
}

var data = []byte(`
<?xml version="1.0" encoding="utf-8"?>
<COLLADA xmlns="http://www.collada.org/2005/11/COLLADASchema" version="1.4.1">
  <asset>
    <contributor>
      <author>Blender User</author>
      <authoring_tool>Blender 2.77.0 commit date:2016-04-05, commit time:18:12, hash:abf6f08</authoring_tool>
    </contributor>
    <created>2016-07-26T16:24:10</created>
    <modified>2016-07-26T16:24:10</modified>
    <unit name="meter" meter="1"/>
    <up_axis>Z_UP</up_axis>
  </asset>
  <library_images/>
  <library_effects>
    <effect id="Material-effect">
      <profile_COMMON>
        <technique sid="common">
          <phong>
            <emission>
              <color sid="emission">0 0 0 1</color>
            </emission>
            <ambient>
              <color sid="ambient">0 0 0 1</color>
            </ambient>
            <diffuse>
              <color sid="diffuse">0.64 0.64 0.64 1</color>
            </diffuse>
            <specular>
              <color sid="specular">0.5 0.5 0.5 1</color>
            </specular>
            <shininess>
              <float sid="shininess">50</float>
            </shininess>
            <index_of_refraction>
              <float sid="index_of_refraction">1</float>
            </index_of_refraction>
          </phong>
        </technique>
      </profile_COMMON>
    </effect>
  </library_effects>
  <library_materials>
    <material id="Material-material" name="Material">
      <instance_effect url="#Material-effect"/>
    </material>
  </library_materials>
  <library_geometries>
    <geometry id="Cube-mesh" name="Cube">
      <mesh>
        <source id="Cube-mesh-positions">
          <float_array id="Cube-mesh-positions-array" count="24">1 1 -1 1 -1 -1 -1 -0.9999998 -1 -0.9999997 1 -1 1 0.9999995 1 0.9999994 -1.000001 1 -1 -0.9999997 1 -1 1 1</float_array>
          <technique_common>
            <accessor source="#Cube-mesh-positions-array" count="8" stride="3">
              <param name="X" type="float"/>
              <param name="Y" type="float"/>
              <param name="Z" type="float"/>
            </accessor>
          </technique_common>
        </source>
        <source id="Cube-mesh-normals">
          <float_array id="Cube-mesh-normals-array" count="36">0 0 -1 0 0 1 1 0 -2.38419e-7 0 -1 -4.76837e-7 -1 2.38419e-7 -1.49012e-7 2.68221e-7 1 2.38419e-7 0 0 -1 0 0 1 1 -5.96046e-7 3.27825e-7 -4.76837e-7 -1 0 -1 2.38419e-7 -1.19209e-7 2.08616e-7 1 0</float_array>
          <technique_common>
            <accessor source="#Cube-mesh-normals-array" count="12" stride="3">
              <param name="X" type="float"/>
              <param name="Y" type="float"/>
              <param name="Z" type="float"/>
            </accessor>
          </technique_common>
        </source>
        <vertices id="Cube-mesh-vertices">
          <input semantic="POSITION" source="#Cube-mesh-positions"/>
        </vertices>
        <polylist material="Material-material" count="12">
          <input semantic="VERTEX" source="#Cube-mesh-vertices" offset="0"/>
          <input semantic="NORMAL" source="#Cube-mesh-normals" offset="1"/>
          <vcount>3 3 3 3 3 3 3 3 3 3 3 3 </vcount>
          <p>0 0 2 0 3 0 7 1 5 1 4 1 4 2 1 2 0 2 5 3 2 3 1 3 2 4 7 4 3 4 0 5 7 5 4 5 0 6 1 6 2 6 7 7 6 7 5 7 4 8 5 8 1 8 5 9 6 9 2 9 2 10 6 10 7 10 0 11 3 11 7 11</p>
        </polylist>
      </mesh>
    </geometry>
  </library_geometries>
  <library_controllers/>
  <library_visual_scenes>
    <visual_scene id="Scene" name="Scene">
      <node id="Cube" name="Cube" type="NODE">
        <matrix sid="transform">1 0 0 0 0 1 0 0 0 0 1 1 0 0 0 1</matrix>
        <instance_geometry url="#Cube-mesh" name="Cube">
          <bind_material>
            <technique_common>
              <instance_material symbol="Material-material" target="#Material-material"/>
            </technique_common>
          </bind_material>
        </instance_geometry>
      </node>
    </visual_scene>
  </library_visual_scenes>
  <scene>
    <instance_visual_scene url="#Scene"/>
  </scene>
</COLLADA>
`)
