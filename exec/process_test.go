package exec

import (
	"strings"
	"testing"
)

const TestText = "MIT License\n\nCopyright (c) 2018 Rao-Mengnan\n\n" +
	"Permission is hereby granted, free of charge, to any person obtaining a copy\n" +
	"of this software and associated documentation files (the \"Software\"), to deal\n" +
	"in the Software without restriction, including without limitation the rights\n" +
	"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell\n" +
	"copies of the Software, and to permit persons to whom the Software is\n" +
	"furnished to do so, subject to the following conditions:\n\n" +
	"" +
	"The above copyright notice and this permission notice shall be included in all\n" +
	"copies or substantial portions of the Software.\n\n" +
	"" +
	"THE SOFTWARE IS PROVIDED \"AS IS\", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR\n" +
	"IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,\n" +
	"FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE\n" +
	"AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER\n" +
	"LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,\n" +
	"OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE\n" +
	"SOFTWARE.\n\n"

func TestEncryptAndDecrypt(t *testing.T) {
	origin := TestText
	encrypted, err := encryptText(origin, string(SK))
	if strings.Compare(origin, encrypted) == 0 {
		t.Errorf("Failed to encrypt")
	}
	decrypted, err := decryptText(string(encrypted), string(SK))
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if strings.Compare(decrypted, origin) != 0 {
		t.Errorf("Failed to decrypt")
	}
}

func BenchmarkEncryptAndDecrypt(b *testing.B) {

	for i := 0; i < b.N; i++ {
		key := "asdfaasdadafdafsadsfsafweqfsda"
		encrypted, _ := encryptText(TestText, key)
		decrypted, _ := decryptText(string(encrypted), key)
		if strings.Compare(TestText, decrypted) != 0 {
			b.Errorf("Encrtpt and decrypt got unexpect result.")
		}
	}
}
