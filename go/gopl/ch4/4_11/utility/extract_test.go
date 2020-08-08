package utility

import "testing"

func TestExtractTitleAndBodyWithHappyPath(t *testing.T) {
	if title, body, isBodyExist := ExtractTitleAndBody("title\n\nbody"); title != "title" || body != "body" || !isBodyExist {
		t.Error("failed the happy path")
	}
}

func TestExtractTitleAndBodyWithMultiTitle(t *testing.T) {
	if title, body, isBodyExist := ExtractTitleAndBody("title\nmulti\n\nbody"); title != "title\nmulti" || body != "body" || !isBodyExist {
		t.Error("failed to handle multi line title")
	}
}

func TestExtractTitleAndBodyWithMultiBody(t *testing.T) {
	if title, body, isBodyExist := ExtractTitleAndBody("title\n\nbody\nmulti body"); title != "title" || body != "body\nmulti body" || !isBodyExist {
		t.Error("failed to handle multi line body")
	}

	if title, body, isBodyExist := ExtractTitleAndBody("#a\ntitle\n#b\n\n#c\nbody\n#e\n#f"); title != "title" || body != "body" || !isBodyExist {
		t.Error("failed to ignore comment")
	}

	if title, _, isBodyExist := ExtractTitleAndBody("title"); title != "title" || isBodyExist {
		t.Error("failed to handle empty body")
	}
}

func TestExtractTitleAndBodyWithComment(t *testing.T) {
	if title, body, isBodyExist := ExtractTitleAndBody("#a\ntitle\n#b\n\n#c\nbody\n#e\n#f"); title != "title" || body != "body" || !isBodyExist {
		t.Error("failed to ignore comment")
	}
}

func TestExtractTitleAndBodyWithNullBody(t *testing.T) {
	if title, _, isBodyExist := ExtractTitleAndBody("title"); title != "title" || isBodyExist {
		t.Error("failed to handle null body")
	}
}
