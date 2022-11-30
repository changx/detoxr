package ds

const (
	RR_A          = 1 + iota   //1	a host address	[RFC1035]
	RR_NS                      //2	an authoritative name server	[RFC1035]
	RR_MD                      //3	a mail destination (OBSOLETE - use MX)	[RFC1035]
	RR_MF                      //4	a mail forwarder (OBSOLETE - use MX)	[RFC1035]
	RR_CNAME                   //5	the canonical name for an alias	[RFC1035]
	RR_SOA                     //6	marks the start of a zone of authority	[RFC1035]
	RR_MB                      //7	a mailbox domain name (EXPERIMENTAL)	[RFC1035]
	RR_MG                      //8	a mail group member (EXPERIMENTAL)	[RFC1035]
	RR_MR                      //9	a mail rename domain name (EXPERIMENTAL)	[RFC1035]
	RR_NULL                    //10	a null RR (EXPERIMENTAL)	[RFC1035]
	RR_WKS                     //11	a well known service description	[RFC1035]
	RR_PTR                     //12	a domain name pointer	[RFC1035]
	RR_HINFO                   //13	host information	[RFC1035]
	RR_MINFO                   //14	mailbox or mail list information	[RFC1035]
	RR_MX                      //15	mail exchange	[RFC1035]
	RR_TXT                     //16	text strings	[RFC1035]
	RR_RP                      //17	for Responsible Person	[RFC1183]
	RR_AFSDB                   //18	for AFS Data Base location	[RFC1183][RFC5864]
	RR_X25                     //19	for X.25 PSDN address	[RFC1183]
	RR_ISDN                    //20	for ISDN address	[RFC1183]
	RR_RT                      //21	for Route Through	[RFC1183]
	RR_NSAP                    //22	for NSAP address, NSAP style A record	[RFC1706]
	RR_NSAP_PTR                //23	for domain name pointer, NSAP style	[RFC1706]
	RR_SIG                     //24	for security signature	[RFC2536][RFC2931][RFC3110][RFC4034]
	RR_KEY                     //25	for security key	[RFC2536][RFC2539][RFC3110][RFC4034]
	RR_PX                      //26	X.400 mail mapping information	[RFC2163]
	RR_GPOS                    //27	Geographical Position	[RFC1712]
	RR_AAAA                    //28	IP6 Address	[RFC3596]
	RR_LOC                     //29	Location Information	[RFC1876]
	RR_NXT                     //30	Next Domain (OBSOLETE)	[RFC2535][RFC3755]
	RR_EID                     //31	Endpoint Identifier	[Michael_Patton][http://ana-3.lcs.mit.edu/~jnc/nimrod/dns.txt]	 //1995-06
	RR_NIMLOC                  //32	Nimrod Locator	[1][Michael_Patton][http://ana-3.lcs.mit.edu/~jnc/nimrod/dns.txt]	 //1995-06
	RR_SRV                     //33	Server Selection	[1][RFC2782]
	RR_ATMA                    //34	ATM Address	[ ATM Forum Technical Committee, "ATM Name System, V2.0", Doc ID: AF-DANS-0152.000, July //2000. Available from and held in escrow by IANA.]
	RR_NAPTR                   //35	Naming Authority Pointer	[RFC3403]
	RR_KX                      //36	Key Exchanger	[RFC2230]
	RR_CERT                    //37	CERT	[RFC4398]
	RR_A6                      //38	A6 (OBSOLETE - use AAAA)	[RFC2874][RFC3226][RFC6563]
	RR_DNAME                   //39	DNAME	[RFC6672]
	RR_SINK                    //40	SINK	[Donald_E_Eastlake][draft-eastlake-kitchen-sink]	 //1997-11
	RR_OPT                     //41	OPT	[RFC3225][RFC6891]
	RR_APL                     //42	APL	[RFC3123]
	RR_DS                      //43	Delegation Signer	[RFC4034]
	RR_SSHFP                   //44	SSH Key Fingerprint	[RFC4255]
	RR_IPSECKEY                //45	IPSECKEY	[RFC4025]
	RR_RRSIG                   //46	RRSIG	[RFC4034]
	RR_NSEC                    //47	NSEC	[RFC4034][RFC9077]
	RR_DNSKEY                  //48	DNSKEY	[RFC4034]
	RR_DHCID                   //49	DHCID	[RFC4701]
	RR_NSEC3                   //50	NSEC3	[RFC5155][RFC9077]
	RR_NSEC3PARAM              //51	NSEC3PARAM	[RFC5155]
	RR_TLSA                    //52	TLSA	[RFC6698]
	RR_SMIMEA                  //53	S/MIME cert association	[RFC8162]	SMIMEA/smimea-completed-template //2015-12-01
	RR_HIP                     //55	Host Identity Protocol	[RFC8005]
	RR_NINFO                   //56	NINFO	[Jim_Reid]	NINFO/ninfo-completed-template //2008-01-21
	RR_RKEY                    //57	RKEY	[Jim_Reid]	RKEY/rkey-completed-template //2008-01-21
	RR_TALINK                  //58	Trust Anchor LINK	[Wouter_Wijngaards]	TALINK/talink-completed-template //2010-02-17
	RR_CDS                     //59	Child DS	[RFC7344]	CDS/cds-completed-template //2011-06-06
	RR_CDNSKEY                 //60	DNSKEY(s) the Child wants reflected in DS	[RFC7344]	 //2014-06-16
	RR_OPENPGPKEY              //61	OpenPGP Key	[RFC7929]	OPENPGPKEY/openpgpkey-completed-template //2014-08-12
	RR_CSYNC                   //62	Child-To-Parent Synchronization	[RFC7477]	 //2015-01-27
	RR_ZONEMD                  //63	Message Digest Over Zone Data	[RFC8976]	ZONEMD/zonemd-completed-template //2018-12-12
	RR_SVCB                    //64	General Purpose Service Binding	[RFC-ietf-dnsop-svcb-https-10]	SVCB/svcb-completed-template //2020-06-30
	RR_HTTPS                   //65	Service Binding type for use with HTTP	[RFC-ietf-dnsop-svcb-https-10]	HTTPS/https-completed-template //2020-06-30
	RR_SPF        = 35 + iota  //99		[RFC7208]
	RR_UINFO                   //100		[IANA-Reserved]
	RR_UID                     //101		[IANA-Reserved]
	RR_GID                     //102		[IANA-Reserved]
	RR_UNSPEC                  //103		[IANA-Reserved]
	RR_NID                     //104		[RFC6742]	ILNP/nid-completed-template
	RR_L32                     //105		[RFC6742]	ILNP/l32-completed-template
	RR_L64                     //106		[RFC6742]	ILNP/l64-completed-template
	RR_LP                      //107		[RFC6742]	ILNP/lp-completed-template
	RR_EUI48                   //108	an EUI-48 address	[RFC7043]	EUI48/eui48-completed-template //2013-03-27
	RR_EUI64                   //109	an EUI-64 address	[RFC7043]	EUI64/eui64-completed-template //2013-03-27
	RR_TKEY       = 174 + iota //249	Transaction Key	[RFC2930]
	RR_TSIG                    //250	Transaction Signature	[RFC8945]
	RR_IXFR                    //251	incremental transfer	[RFC1995]
	RR_AXFR                    //252	transfer of an entire zone	[RFC1035][RFC5936]
	RR_MAILB                   //253	mailbox-related RRs (MB, MG or MR)	[RFC1035]
	RR_MAILA                   //254	mail agent RRs (OBSOLETE - see MX)	[RFC1035]
	RR_ANY                     //255	A request for some or all records the server has available	[RFC1035][RFC6895][RFC8482]
	RR_URI                     //256	URI	[RFC7553]	URI/uri-completed-template //2011-02-22
	RR_CAA                     //257	Certification Authority Restriction	[RFC8659]	CAA/caa-completed-template //2011-04-07
	RR_AVC                     //258	Application Visibility and Control	[Wolfgang_Riedel]	AVC/avc-completed-template //2016-02-26
	RR_DOA                     //259	Digital Object Architecture	[draft-durand-doa-over-dns]	DOA/doa-completed-template //2017-08-30
	RR_AMTRELAY                //260	Automatic Multicast Tunneling Relay	[RFC8777]
)
