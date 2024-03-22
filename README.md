# Resume Styler

## Description

This is a simple tool built using Go that takes in a resume.md and styles it outputing a pdf

## How to use

1. Create a markdown file with your resume
2. Create a custom style in css (If run without a custom style, the default) style will be used
3. Try before building

```bash
go run main.go Resume-Dev.md
# or
go run main.go Resume-Dev.md customStyle.css
```

## Resume Structure

Styling the resume is easiest when a clear hierarchy is established. The following is a basic pattern for a resume:

```markdown
# Name (h1) {#my-custom-id}

## Section Heading (h2)

### Subsection Heading (h3)

--- Section Line (Horizontal Rule)

- Section Content List (Unordered List)

> Section Content (Blockquote)
```

## Custom Styles

The custom style should be a css file that will be used to style the resume. The following is a basic structure for a custom style:

```css
body {
  font-family:
    FiraCode Nerd Font,
    Arial,
    sans-serif;
  padding: 20px;
}

#my-custom-id {
  color: gray;
}

h1 {
  font-size: 17px;
  margin-bottom: 2px;
}

p {
  font-size: 12px;
  line-height: 1.5;
  padding: 10px;
}

blockquote {
  text-align: left;
  line-height: 1.6;
}

h2 {
  font-size: 18px;
  margin-bottom: 5px;
}

h3 {
  font-size: 14px;
  margin-bottom: 5px;
}

p :not(blockquote p) {
  text-align: center;
  font-size: 12px;
  line-height: 1.5;
}

ul {
  list-style-type: square; /* Remove default bullets */
  padding: 0;
  margin: 0;
}

ul li {
  font-size: 12px;
  padding: 5px 0; /* Add some padding to each list item */
  margin-left: 40px;
}

ul li:last-child {
  border-bottom: none; /* Remove the border from the last list item */
}
```
