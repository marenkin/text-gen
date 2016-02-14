package textgen

type Status int

const (
	StatusEnabled Status = iota
	StatusDisabled
	StatusTransaction
)

type IMetatype interface {
	Id() int
	Name() string

	LinkedMetatypes() []IMetatype
	ReverseLinkedMetatype() []IMetatype

	Status() Status
}

type IMetatypeLink interface {
	Id() int

	Metatype() IMetatype

	Status() Status
}

type IMetaparam interface {
	Id() int // or index in case of in-memory dataset
	Name() string
	Metatype() IMetatype
	LanguageCode() string

	Status() Status
}

type IContent interface {
	Id() int
	Ref() string
	Metatype() IMetatype
	GetParams(metaparamId int) []IContentparam

	Status() Status
}

type IContentparam interface {
	Metaparam() IMetaparam
	Value() string

	Status() Status
}

type IContentSource interface {
	// add functions to populate content
	// it helps to implement content source conversion function
	//  or content source comparison ( matching function)

	AddMetatype(name string) (id int)
	AddMetaparam(metatypeId int, name string) (id int)
	AddMetatypeLink(metatypeId int, linkedMetatypeId int)

	AddContent(metatypeId int, ref string) (id int)
	AddContentParam(contentId int, metatypeParamId int, value string) (id int)
	AddContentLink(contentId int, linkedContentId int)

	Init() (err error)

	TransactionBegin()  // prepare structures etc
	TransactionCommit() // finalise updates

	ListContent(metatypeId int, lastContentId int, pageSize int) []IContent
	ListLinkedContent(parentContentId int, linkedMetatypeId int, lastLinkedContentId int, pageSize int) []IContent
}

// content manager to be move to a separate file

type ContentManager struct {
	ContentSource IContentSource

	// copy content source

	// compare two content sources

	// clone datasources :  create a new and copy

}

type IContentSourceReader interface {
	ListContent(lastContentId int, pageSize int) []IContent
}

type ContentReader struct {
	ContentSourceReader IContentSourceReader
	Items               []IContent
	CurrentIndex        int
}

func (m *ContentReader) Reset() {
	m.CurrentIndex = -1
}

func (m *ContentReader) Read() bool {

	m.CurrentIndex++

	return !m.EOF()
}

func (m *ContentReader) Current() IContent {
	return m.Items[m.CurrentIndex]
}

func (m *ContentReader) EOF() bool {
	return len(m.Items) <= m.CurrentIndex
}

type Step struct {
	Direction   bool
	Metatype    IMetatype
	RelMetatype IMetatype
}

// to be moved to separate file ( in-memory content source )
type Metatype struct {
	Id int // index

	Name       string
	Metaparams []*Metaparam

	Content []*Content // index is an Id

	LinkedMetatypes map[string]int // key is a linked metatype name, value is index

	Reverselinkedmetatypes map[string]int // metatypes referencing currenct metatype in the direct link

	paths [][]*Step // []Step is indexed by Metatype->Index
}

type Metaparam struct {
	Id int // is param index

	Metatype     *Metatype
	Name         string
	LanguageCode string
}

type Content struct {
	Id int // content index

	Metatype *Metatype
	Ref      string
	Params   [][]*Contentparam //*Contentparam //indexed by Metaparam->index
}

type Contentparam struct {
	Status int // change to enum
	Value  string
}
