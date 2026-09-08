//nolint
package sanitize

func init() {
	serviceFields[Response]["oks"] = map[string]string{
		"kubeconfig": Sensitive,
	}
}