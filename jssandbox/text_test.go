package jssandbox

import (
	"context"
	"testing"
)

func TestReplaceText(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = replaceText("Hello World", "World", "Go");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("replaceText() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("replaceText()应该返回success: true")
	}

	resultText := resultObj.Get("result").String()
	if resultText != "Hello Go" {
		t.Errorf("replaceText()结果 = %v, 期望 = Hello Go", resultText)
	}
}

func TestReplaceTextAll(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = replaceText("foo bar foo", "foo", "baz", true);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("replaceText() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	resultText := resultObj.Get("result").String()
	if resultText != "baz bar baz" {
		t.Errorf("replaceText()结果 = %v, 期望 = baz bar baz", resultText)
	}
}

func TestReplaceRegex(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = replaceRegex("Hello 123 World", "\\d+", "XXX");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("replaceRegex() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("replaceRegex()应该返回success: true")
	}

	resultText := resultObj.Get("result").String()
	if resultText != "Hello XXX World" {
		t.Errorf("replaceRegex()结果 = %v, 期望 = Hello XXX World", resultText)
	}
}

func TestMatchMarkdownHeaders(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	markdown := `# Title 1
Some content
## Title 2
More content
### Title 3
`

	code := `
		var result = matchMarkdownHeaders(` + "`" + markdown + "`" + `);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("matchMarkdownHeaders() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("matchMarkdownHeaders()应该返回success: true")
	}

	count := resultObj.Get("count").ToInteger()
	if count != 3 {
		t.Errorf("matchMarkdownHeaders() count = %v, 期望 = 3", count)
	}

	headers := resultObj.Get("headers")
	headersObj := headers.ToObject(sb.vm)
	firstHeader := headersObj.Get("0").ToObject(sb.vm)
	firstLevel := firstHeader.Get("level").ToInteger()
	if firstLevel != 1 {
		t.Errorf("第一个标题级别 = %v, 期望 = 1", firstLevel)
	}
}

func TestMatchMarkdownHeaderByLevel(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	markdown := `# Title 1
## Title 2
## Title 2-2
### Title 3
`

	code := `
		var result = matchMarkdownHeaderByLevel(` + "`" + markdown + "`" + `, 2);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("matchMarkdownHeaderByLevel() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	count := resultObj.Get("count").ToInteger()
	if count != 2 {
		t.Errorf("matchMarkdownHeaderByLevel() count = %v, 期望 = 2", count)
	}
}

func TestMatchMarkdownImages(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	markdown := `![Alt text](https://example.com/image.png)
![Another image](https://example.com/image2.jpg "Image Title")
Some text
`

	code := `
		var result = matchMarkdownImages(` + "`" + markdown + "`" + `);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("matchMarkdownImages() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("matchMarkdownImages()应该返回success: true")
	}

	count := resultObj.Get("count").ToInteger()
	if count != 2 {
		t.Errorf("matchMarkdownImages() count = %v, 期望 = 2", count)
	}

	images := resultObj.Get("images")
	imagesObj := images.ToObject(sb.vm)
	firstImage := imagesObj.Get("0").ToObject(sb.vm)
	firstUrl := firstImage.Get("url").String()
	if firstUrl != "https://example.com/image.png" {
		t.Errorf("第一个图片URL = %v, 期望 = https://example.com/image.png", firstUrl)
	}

	secondImage := imagesObj.Get("1").ToObject(sb.vm)
	secondTitle := secondImage.Get("title").String()
	if secondTitle != "Image Title" {
		t.Errorf("第二个图片标题 = %v, 期望 = Image Title", secondTitle)
	}
}

func TestMatchMarkdownCodeBlocks(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	markdown := "Here is some code:\n```go\nfunc main() {\n    println(\"Hello\")\n}\n```\nAnd inline code: `code`\n"

	code := `
		var result = matchMarkdownCodeBlocks(` + "`" + markdown + "`" + `, true);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("matchMarkdownCodeBlocks() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("matchMarkdownCodeBlocks()应该返回success: true")
	}

	count := resultObj.Get("count").ToInteger()
	if count != 2 {
		t.Errorf("matchMarkdownCodeBlocks() count = %v, 期望 = 2", count)
	}

	codeBlocks := resultObj.Get("codeBlocks")
	codeBlocksObj := codeBlocks.ToObject(sb.vm)
	firstBlock := codeBlocksObj.Get("0").ToObject(sb.vm)
	firstType := firstBlock.Get("type").String()
	if firstType != "fenced" {
		t.Errorf("第一个代码块类型 = %v, 期望 = fenced", firstType)
	}

	firstLanguage := firstBlock.Get("language").String()
	if firstLanguage != "go" {
		t.Errorf("第一个代码块语言 = %v, 期望 = go", firstLanguage)
	}
}

func TestReplaceTemplate(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var template = "Hello {{name}}, you are {{age}} years old.";
		var data = {name: "Alice", age: 30};
		var result = replaceTemplate(template, data);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("replaceTemplate() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("replaceTemplate()应该返回success: true")
	}

	resultText := resultObj.Get("result").String()
	expected := "Hello Alice, you are 30 years old."
	if resultText != expected {
		t.Errorf("replaceTemplate()结果 = %v, 期望 = %v", resultText, expected)
	}
}

func TestReplaceTemplateAdvanced(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var template = "Hello {{name|Guest}}, you are {{age|unknown}} years old.";
		var data = {name: "Alice"};
		var result = replaceTemplateAdvanced(template, data);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("replaceTemplateAdvanced() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("replaceTemplateAdvanced()应该返回success: true")
	}

	resultText := resultObj.Get("result").String()
	expected := "Hello Alice, you are unknown years old."
	if resultText != expected {
		t.Errorf("replaceTemplateAdvanced()结果 = %v, 期望 = %v", resultText, expected)
	}
}

func TestExtractMarkdownStructure(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	markdown := `# Chapter 1
Content here
## Section 1.1
More content
## Section 1.2
Even more
### Subsection 1.2.1
Details
# Chapter 2
Final content
`

	code := `
		var result = extractMarkdownStructure(` + "`" + markdown + "`" + `);
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("extractMarkdownStructure() error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	success := resultObj.Get("success")
	if !success.ToBoolean() {
		t.Error("extractMarkdownStructure()应该返回success: true")
	}

	structure := resultObj.Get("structure")
	structureObj := structure.ToObject(sb.vm)
	length := structureObj.Get("length").ToInteger()
	if length != 2 {
		t.Errorf("extractMarkdownStructure() 顶级节点数 = %v, 期望 = 2", length)
	}

	// 检查第一个章节的子节点
	chapter1 := structureObj.Get("0").ToObject(sb.vm)
	chapter1Children := chapter1.Get("children")
	chapter1ChildrenObj := chapter1Children.ToObject(sb.vm)
	childrenLength := chapter1ChildrenObj.Get("length").ToInteger()
	if childrenLength != 2 {
		t.Errorf("Chapter 1 子节点数 = %v, 期望 = 2", childrenLength)
	}
}

func TestReplaceTextEdgeCases(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	// 测试空字符串
	code := `
		var result = replaceText("", "old", "new");
		result;
	`
	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("replaceText() with empty string error = %v", err)
	}
	resultObj := result.ToObject(sb.vm)
	resultText := resultObj.Get("result").String()
	if resultText != "" {
		t.Errorf("空字符串替换结果 = %v, 期望 = 空字符串", resultText)
	}

	// 测试不存在的字符串
	code2 := `
		var result = replaceText("Hello World", "xyz", "abc");
		result;
	`
	result2, err := sb.Run(code2)
	if err != nil {
		t.Fatalf("replaceText() with non-existent string error = %v", err)
	}
	resultObj2 := result2.ToObject(sb.vm)
	resultText2 := resultObj2.Get("result").String()
	if resultText2 != "Hello World" {
		t.Errorf("不存在的字符串替换结果 = %v, 期望 = Hello World", resultText2)
	}
}

func TestMatchMarkdownHeadersEmpty(t *testing.T) {
	ctx := context.Background()
	sb := NewSandbox(ctx)
	defer sb.Close()

	code := `
		var result = matchMarkdownHeaders("");
		result;
	`

	result, err := sb.Run(code)
	if err != nil {
		t.Fatalf("matchMarkdownHeaders() with empty string error = %v", err)
	}

	resultObj := result.ToObject(sb.vm)
	count := resultObj.Get("count").ToInteger()
	if count != 0 {
		t.Errorf("空字符串的标题数 = %v, 期望 = 0", count)
	}
}

