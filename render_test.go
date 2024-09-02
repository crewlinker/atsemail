package atsemail_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"

	"github.com/crewlinker/atsemail"
)

func AssertEmailRender[T proto.Message](tb testing.TB, templateName string, caseIdx int, data T, expf func(g Gomega, txtbuf, htbuf *bytes.Buffer)) {
	tb.Helper()
	g, ctx := NewWithT(tb), context.Background()

	render, err := atsemail.New[T](templateName)
	g.Expect(err).ToNot(HaveOccurred())

	var txtbuf, htbuf bytes.Buffer
	g.Expect(render.Render(&txtbuf, &htbuf, data)).To(Succeed())

	expf(g, &htbuf, &txtbuf)

	SaveScreenshot(ctx, g, fmt.Sprintf("%s_%d", templateName, caseIdx), &htbuf)
}

func SaveScreenshot(ctx context.Context, g Gomega, name string, htbuf *bytes.Buffer) {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var buf []byte
	g.Expect(chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			ft, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to get frame tree: %w", err)
			}

			if err := page.SetDocumentContent(ft.Frame.ID, htbuf.String()).Do(ctx); err != nil {
				return fmt.Errorf("failed to set docment content: %w", err)
			}

			return nil
		}),
		chromedp.FullScreenshot(&buf, 100),
	)).To(Succeed())

	g.Expect(os.WriteFile(filepath.Join("screenshots", name+".png"), buf, 0o644)).To(Succeed())
}
