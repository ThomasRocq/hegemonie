package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

type Vertex interface {
	Id() string

	IsCenter() bool
	SetCenter(v bool)

	IsAnchor() bool
	SetAnchor(v bool)
}

type Edge interface {
	S() Vertex
	D() Vertex
}

type Graph interface {
	AddVertex(e Vertex)
	AddEdge(v Edge)

	GetVertex() []Vertex
	GetEdges() []Edge

	GetAnchors() []Vertex

	Draw(o io.Writer)
}

type memGraph struct {
	vertices []Vertex
	edges    []Edge
}

type memVertex struct {
	id     uint64
	center bool
	anchor bool
}

type memEdge struct {
	s Vertex
	d Vertex
}

type clusteredGraph struct {
	memGraph
	clusters []Graph
}

var nextId uint64

func getNextId() uint64 {
	return atomic.AddUint64(&nextId, 1)
}

func MakeEdge(s, d Vertex) Edge { return &memEdge{s: s, d: d} }
func (e *memEdge) S() Vertex    { return e.s }
func (e *memEdge) D() Vertex    { return e.d }

func MakeVertex() Vertex              { return &memVertex{id: getNextId()} }
func (v *memVertex) Id() string       { return strconv.FormatUint(v.id, 16) }
func (v *memVertex) IsAnchor() bool   { return v.anchor }
func (v *memVertex) IsCenter() bool   { return v.center }
func (v *memVertex) SetAnchor(b bool) { v.anchor = b }
func (v *memVertex) SetCenter(b bool) { v.center = b }

func (g *memGraph) AddEdge(e Edge)      { g.edges = append(g.edges, e) }
func (g *memGraph) AddVertex(v Vertex)  { g.vertices = append(g.vertices, v) }
func (g *memGraph) GetEdges() []Edge    { return g.edges[:] }
func (g *memGraph) GetVertex() []Vertex { return g.vertices[:] }
func (g *memGraph) GetAnchors() []Vertex {
	rc := make([]Vertex, 0)
	for _, v := range g.vertices {
		if v.IsAnchor() {
			rc = append(rc, v)
		}
	}
	return rc
}

func (g *memGraph) Draw(o io.Writer) {
	for _, e := range g.edges {
		fmt.Fprintf(o, "\"%s\" -- \"%s\";\n", e.S().Id(), e.D().Id())
	}
}

func (cg *clusteredGraph) Add(g Graph)        { cg.clusters = append(cg.clusters, g) }
func (cg *clusteredGraph) AddEdge(e Edge)     { cg.edges = append(cg.edges, e) }
func (cg *clusteredGraph) AddVertex(v Vertex) { cg.vertices = append(cg.vertices, v) }
func (cg *clusteredGraph) GetEdges() []Edge {
	rc := make([]Edge, 0)
	for _, g := range cg.clusters {
		rc = append(rc, g.GetEdges()...)
	}
	return rc
}

func (cg *clusteredGraph) GetVertex() []Vertex {
	rc := make([]Vertex, 0)
	for _, g := range cg.clusters {
		rc = append(rc, g.GetVertex()...)
	}
	return rc
}

func (cg *clusteredGraph) GetAnchors() []Vertex {
	rc := make([]Vertex, 0)
	for _, g := range cg.clusters {
		rc = append(rc, g.GetAnchors()...)
	}
	return rc
}

func (cg *clusteredGraph) Draw(o io.Writer) {
	for _, g := range cg.clusters {
		g.Draw(o)
	}

	for _, e := range cg.edges {
		fmt.Fprintf(o, "\"%s\" -- \"%s\";\n", e.S().Id(), e.D().Id())
	}
}

func Cluster(g ...Graph) Graph {
	rc := &clusteredGraph{}
	rc.clusters = append(rc.clusters, g...)
	return rc
}

func genCircle(n uint) Graph {
	if n < 3 {
		panic("A line with length < 2 makes no sense")
	}

	vertices := make([]Vertex, 0)
	for i := uint(0); i < n; i++ {
		vertices = append(vertices, MakeVertex())
	}
	for i, v := range vertices {
		if (i % 2) == 0 {
			v.SetAnchor(true)
		}
	}

	edges := make([]Edge, 0)
	for i := uint(0); i < n; i++ {
		e := MakeEdge(vertices[i], vertices[int(i+1)%len(vertices)])
		edges = append(edges, e)
	}

	return &memGraph{vertices: vertices, edges: edges}
}

func genLine(n uint) Graph {
	if n < 2 {
		panic("A line with length < 2 makes no sense")
	}

	vertices := make([]Vertex, 0)
	for i := uint(0); i < n; i++ {
		vertices = append(vertices, MakeVertex())
	}
	for i, v := range vertices {
		if (i % 2) == 0 {
			v.SetAnchor(true)
		}
	}
	vertices[n-1].SetAnchor(true)

	edges := make([]Edge, 0)
	for i := uint(1); i < n; i++ {
		e := MakeEdge(vertices[i-1], vertices[i])
		edges = append(edges, e)
	}

	return &memGraph{vertices: vertices, edges: edges}
}

func genStar(nb, radius uint) Graph {
	if nb < 3 {
		panic("Star with less than 3 branches makes no sense")
	}
	if radius < 1 {
		panic("Star with radius zero makes no sense")
	}

	clusters := make([]Graph, 0)
	for i := uint(0); i < nb; i++ {
		clusters = append(clusters, genLine(radius))
	}

	return GlueStar(clusters...)
}

func GlueStar(gv ...Graph) Graph {
	if len(gv) < 2 {
		panic("Glueing less than 2 Graphs makes no sense")
	}

	cluster := Cluster(gv...)
	a := MakeVertex()
	cluster.AddVertex(a)
	for _, g := range gv {
		b := PeekAnchor(g)
		cluster.AddEdge(MakeEdge(a, b))
		b.SetAnchor(false)
	}
	a.SetAnchor(false)
	return cluster
}

func GlueChain(gv ...Graph) Graph {
	if len(gv) < 2 {
		panic("Glueing less than 2 Graphs makes no sense")
	}

	rc := Cluster(gv[0])
	for _, g := range gv[1:] {
		a := PeekAnchor(rc)
		a.SetAnchor(false)
		b := PeekAnchor(g)
		b.SetAnchor(false)
		rc.AddEdge(MakeEdge(a, b))
		rc.(*clusteredGraph).Add(g)
	}
	return rc
}

func Loop(g Graph) {
	a := PeekAnchor(g)
	a.SetAnchor(false)
	b := PeekAnchor(g)
	b.SetAnchor(false)
	g.AddEdge(MakeEdge(a, b))
}

func PeekAnchor(g Graph) Vertex {
	anchors := g.GetAnchors()
	return anchors[rand.Intn(len(anchors))]
}

func main() {
	rand.Seed(time.Now().UnixNano())
	blocks := []Graph{
		genStar(5, 2),
		genCircle(9),
		genCircle(7),
		genCircle(7),
		genCircle(7),
		genCircle(7),
		genCircle(7),
		genStar(4, 2),
		genStar(3, 2),
	}

	g := GlueChain(blocks...)
	Loop(g)
	Loop(g)

	fmt.Fprint(os.Stdout, "graph g {\n")
	for _, v := range g.GetAnchors() {
		if v.IsAnchor() {
			fmt.Fprintf(os.Stdout, "\"%s\" [color=red];", v.Id())
		}
	}
	g.Draw(os.Stdout)
	fmt.Fprint(os.Stdout, "}\n")
}
