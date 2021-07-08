package Deece

type NoNSType struct{}
type NoENS struct{ domain string }
type NoDNS struct{ domain string }
type CIDmissing struct{ cid string }
type NoNSresolve struct {
	cid    string
	domain string
}
type nopdf struct{ cid string }
type Noipns struct{ ipns string }
type pdfreadfail struct{ name string }
type existCheckFail struct{ name string }
type incorrrectInput struct{}

func (k *NoNSType) Error() string {
	return "Addressing convention not supported"
}
func (m *CIDmissing) Error() string {
	return "CID '" + m.cid + "' could not be resolved."
}
func (l *NoDNS) Error() string {
	return "No DNSLink entry for '" + l.domain + "'."
}
func (n *NoENS) Error() string {
	return "No ENS entry for '" + n.domain + "'."
}
func (o *NoNSresolve) Error() string {
	return "CID '" + o.cid + "' for domain '" + o.domain + "' could not be resolved."
}
func (z *nopdf) Error() string {
	return "File for '" + z.cid + "' is not a pdf."
}
func (d *Noipns) Error() string {
	return "No IPNS entry for '" + d.ipns + "'."
}
func (x *pdfreadfail) Error() string {
	return "Failed to extract data for '" + x.name + "'."
}
func (xx *existCheckFail) Error() string {
	return "Failed to check if '" + xx.name + "' exists."
}
func (zz *incorrrectInput) Error() string {
	return "Input type is not recognised."
}
