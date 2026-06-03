interface TiptapMark {
  type: string;
}

interface TiptapNode {
  type: string;
  attrs?: Record<string, unknown>;
  content?: TiptapNode[];
  text?: string;
  marks?: TiptapMark[];
}

function escapeHtml(str: string): string {
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#x27;");
}

function renderNode(node: TiptapNode): string {
  const children = () => (node.content ?? []).map(renderNode).join("");
  switch (node.type) {
    case "doc":
      return children();
    case "paragraph":
      return `<p>${children()}</p>`;
    case "heading": {
      const level = (node.attrs?.level as number) ?? 1;
      return `<h${level}>${children()}</h${level}>`;
    }
    case "text": {
      let text = escapeHtml(node.text ?? "");
      for (const mark of node.marks ?? []) {
        if (mark.type === "bold") text = `<strong>${text}</strong>`;
        else if (mark.type === "italic") text = `<em>${text}</em>`;
      }
      return text;
    }
    default:
      return children();
  }
}

export function tiptapJsonToHtml(doc: TiptapNode): string {
  return renderNode(doc);
}
