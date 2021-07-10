package Deece

//defining the errors used to return throughout functions in the package
type noNSType struct{}
type noENS struct{ domain string }
type noDNS struct{ domain string }
type cIDmissing struct{ cid string }
type noNSresolve struct {
	cid    string
	domain string
}
type nopdf struct{ cid string }
type noipns struct{ ipns string }
type pdfreadfail struct{ name string }
type existCheckFail struct{ name string }
type IncorrrectInput struct{}
type noIndexAdd struct{}

func (k *noNSType) Error() string {
	return "Addressing convention not supported"
}
func (m *cIDmissing) Error() string {
	return "CID '" + m.cid + "' could not be resolved."
}
func (l *noDNS) Error() string {
	return "No DNSLink entry for '" + l.domain + "'."
}
func (n *noENS) Error() string {
	return "No ENS entry for '" + n.domain + "'."
}
func (o *noNSresolve) Error() string {
	return "CID '" + o.cid + "' for domain '" + o.domain + "' could not be resolved."
}
func (z *nopdf) Error() string {
	return "File for '" + z.cid + "' is not a pdf."
}
func (d *noipns) Error() string {
	return "No IPNS entry for '" + d.ipns + "'."
}
func (x *pdfreadfail) Error() string {
	return "Failed to extract data for '" + x.name + "'."
}
func (xx *existCheckFail) Error() string {
	return "Failed to check if '" + xx.name + "' exists."
}
func (zz *IncorrrectInput) Error() string {
	return "Input type is not recognised."
}

func (zz *noIndexAdd) Error() string {
	return "Not able to add to the index."
}
