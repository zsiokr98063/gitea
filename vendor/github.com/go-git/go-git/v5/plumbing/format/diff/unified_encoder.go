	"bytes"
const (
	diffInit = "diff --git a/%s b/%s\n"
	chunkStart  = "@@ -"
	chunkMiddle = " +"
	chunkEnd    = " @@%s\n"
	chunkCount  = "%d,%d"
	noFilePath = "/dev/null"
	aDir       = "a/"
	bDir       = "b/"
	fPath  = "--- %s\n"
	tPath  = "+++ %s\n"
	binary = "Binary files %s and %s differ\n"

	addLine    = "+%s%s"
	deleteLine = "-%s%s"
	equalLine  = " %s%s"
	noNewLine  = "\n\\ No newline at end of file\n"

	oldMode         = "old mode %o\n"
	newMode         = "new mode %o\n"
	deletedFileMode = "deleted file mode %o\n"
	newFileMode     = "new file mode %o\n"

	renameFrom     = "from"
	renameTo       = "to"
	renameFileMode = "rename %s %s\n"

	indexAndMode = "index %s..%s %o\n"
	indexNoMode  = "index %s..%s\n"

	DefaultContextLines = 3
// UnifiedEncoder encodes an unified diff into the provided Writer.
// There are some unsupported features:
//     - Similarity index for renames
//     - Sort hash representation
	// ctxLines is the count of unchanged lines that will appear
	// surrounding a change.
	ctxLines int
	buf bytes.Buffer
func NewUnifiedEncoder(w io.Writer, ctxLines int) *UnifiedEncoder {
	return &UnifiedEncoder{ctxLines: ctxLines, Writer: w}
	e.printMessage(patch.Message())
	if err := e.encodeFilePatch(patch.FilePatches()); err != nil {
		return err
	_, err := e.buf.WriteTo(e)

	return err
}

func (e *UnifiedEncoder) encodeFilePatch(filePatches []FilePatch) error {
	for _, p := range filePatches {
		f, t := p.Files()
		if err := e.header(f, t, p.IsBinary()); err != nil {
			return err
		}

		g := newHunksGenerator(p.Chunks(), e.ctxLines)
		for _, c := range g.Generate() {
			c.WriteTo(&e.buf)
	return nil
func (e *UnifiedEncoder) printMessage(message string) {
	isEmpty := message == ""
	hasSuffix := strings.HasSuffix(message, "\n")
	if !isEmpty && !hasSuffix {
		message += "\n"
	e.buf.WriteString(message)
}

func (e *UnifiedEncoder) header(from, to File, isBinary bool) error {
	case from == nil && to == nil:
		return nil

		fmt.Fprintf(&e.buf, diffInit, from.Path(), to.Path())

			fmt.Fprintf(&e.buf, oldMode+newMode, from.Mode(), to.Mode())

			fmt.Fprintf(&e.buf,
				renameFileMode+renameFileMode,
				renameFrom, from.Path(), renameTo, to.Path())

			fmt.Fprintf(&e.buf, indexNoMode, from.Hash(), to.Hash())
			fmt.Fprintf(&e.buf, indexAndMode, from.Hash(), to.Hash(), from.Mode())

			e.pathLines(isBinary, aDir+from.Path(), bDir+to.Path())
		fmt.Fprintf(&e.buf, diffInit, to.Path(), to.Path())
		fmt.Fprintf(&e.buf, newFileMode, to.Mode())
		fmt.Fprintf(&e.buf, indexNoMode, plumbing.ZeroHash, to.Hash())
		e.pathLines(isBinary, noFilePath, bDir+to.Path())
		fmt.Fprintf(&e.buf, diffInit, from.Path(), from.Path())
		fmt.Fprintf(&e.buf, deletedFileMode, from.Mode())
		fmt.Fprintf(&e.buf, indexNoMode, from.Hash(), plumbing.ZeroHash)
		e.pathLines(isBinary, aDir+from.Path(), noFilePath)
	return nil
func (e *UnifiedEncoder) pathLines(isBinary bool, fromPath, toPath string) {
	format := fPath + tPath
		format = binary

	fmt.Fprintf(&e.buf, format, fromPath, toPath)
func (c *hunksGenerator) Generate() []*hunk {
	for i, chunk := range c.chunks {
		ls := splitLines(chunk.Content())
		lsLen := len(ls)
			c.fromLine += lsLen
			c.toLine += lsLen
			c.processEqualsLines(ls, i)
			if lsLen != 0 {
				c.fromLine++
			c.processHunk(i, chunk.Type())
			c.fromLine += lsLen - 1
			c.current.AddOp(chunk.Type(), ls...)
			if lsLen != 0 {
				c.toLine++
			c.processHunk(i, chunk.Type())
			c.toLine += lsLen - 1
			c.current.AddOp(chunk.Type(), ls...)
		if i == len(c.chunks)-1 && c.current != nil {
			c.hunks = append(c.hunks, c.current)
	return c.hunks
func (c *hunksGenerator) processHunk(i int, op Operation) {
	if c.current != nil {
	linesBefore := len(c.beforeContext)
	if linesBefore > c.ctxLines {
		ctxPrefix = " " + c.beforeContext[linesBefore-c.ctxLines-1]
		c.beforeContext = c.beforeContext[linesBefore-c.ctxLines:]
		linesBefore = c.ctxLines
	c.current = &hunk{ctxPrefix: strings.TrimSuffix(ctxPrefix, "\n")}
	c.current.AddOp(Equal, c.beforeContext...)
		c.current.fromLine, c.current.toLine =
			c.addLineNumbers(c.fromLine, c.toLine, linesBefore, i, Add)
		c.current.toLine, c.current.fromLine =
			c.addLineNumbers(c.toLine, c.fromLine, linesBefore, i, Delete)
	c.beforeContext = nil
// addLineNumbers obtains the line numbers in a new chunk
func (c *hunksGenerator) addLineNumbers(la, lb int, linesBefore int, i int, op Operation) (cla, clb int) {
	case linesBefore != 0 && c.ctxLines != 0:
		if lb > c.ctxLines {
			clb = lb - c.ctxLines + 1
	case c.ctxLines == 0:
	case i != len(c.chunks)-1:
		next := c.chunks[i+1]
func (c *hunksGenerator) processEqualsLines(ls []string, i int) {
	if c.current == nil {
		c.beforeContext = append(c.beforeContext, ls...)
	c.afterContext = append(c.afterContext, ls...)
	if len(c.afterContext) <= c.ctxLines*2 && i != len(c.chunks)-1 {
		c.current.AddOp(Equal, c.afterContext...)
		c.afterContext = nil
		ctxLines := c.ctxLines
		if ctxLines > len(c.afterContext) {
			ctxLines = len(c.afterContext)
		c.current.AddOp(Equal, c.afterContext[:ctxLines]...)
		c.hunks = append(c.hunks, c.current)
		c.current = nil
		c.beforeContext = c.afterContext[ctxLines:]
		c.afterContext = nil
var splitLinesRE = regexp.MustCompile(`[^\n]*(\n|$)`)

	out := splitLinesRE.FindAllString(s, -1)
func (c *hunk) WriteTo(buf *bytes.Buffer) {
	buf.WriteString(chunkStart)
	if c.fromCount == 1 {
		fmt.Fprintf(buf, "%d", c.fromLine)
		fmt.Fprintf(buf, chunkCount, c.fromLine, c.fromCount)
	buf.WriteString(chunkMiddle)
	if c.toCount == 1 {
		fmt.Fprintf(buf, "%d", c.toLine)
		fmt.Fprintf(buf, chunkCount, c.toLine, c.toCount)
	fmt.Fprintf(buf, chunkEnd, c.ctxPrefix)
	for _, d := range c.ops {
		buf.WriteString(d.String())
func (c *hunk) AddOp(t Operation, s ...string) {
	ls := len(s)
		c.toCount += ls
		c.fromCount += ls
		c.toCount += ls
		c.fromCount += ls
	for _, l := range s {
		c.ops = append(c.ops, &op{l, t})
func (o *op) String() string {
	var prefix, suffix string
	switch o.t {
	case Add:
		prefix = addLine
	case Delete:
		prefix = deleteLine
	case Equal:
		prefix = equalLine
	n := len(o.text)
	if n > 0 && o.text[n-1] != '\n' {
		suffix = noNewLine
	}

	return fmt.Sprintf(prefix, o.text, suffix)