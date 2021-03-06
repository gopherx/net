package nat

var (
	RFC5769SampleRequestPwd = []byte("VOkJxbRl1RmTxUk/WvJxBt")

	RFC5769SampleRequest = []byte{
		0x00, 0x01, 0x00, 0x58, //   Request type and message length
		0x21, 0x12, 0xa4, 0x42, //   Magic cookie
		0xb7, 0xe7, 0xa7, 0x01, //}
		0xbc, 0x34, 0xd6, 0x86, //}  Transaction ID
		0xfa, 0x87, 0xdf, 0xae, //}
		0x80, 0x22, 0x00, 0x10, //   SOFTWARE attribute header
		0x53, 0x54, 0x55, 0x4e, //}
		0x20, 0x74, 0x65, 0x73, //}  User-agent...
		0x74, 0x20, 0x63, 0x6c, //}  ...name
		0x69, 0x65, 0x6e, 0x74, //}
		0x00, 0x24, 0x00, 0x04, //   PRIORITY attribute header
		0x6e, 0x00, 0x01, 0xff, //   ICE priority value
		0x80, 0x29, 0x00, 0x08, //   ICE-CONTROLLED attribute header
		0x93, 0x2f, 0xf9, 0xb1, //}  Pseudo-random tie breaker...
		0x51, 0x26, 0x3b, 0x36, //}   ...for ICE control
		0x00, 0x06, 0x00, 0x09, //   USERNAME attribute header
		0x65, 0x76, 0x74, 0x6a, //}
		0x3a, 0x68, 0x36, 0x76, //}  Username (9 bytes) and padding (3 bytes)
		0x59, 0x20, 0x20, 0x20, //}
		0x00, 0x08, 0x00, 0x14, //   MESSAGE-INTEGRITY attribute header
		0x9a, 0xea, 0xa7, 0x0c, //}
		0xbf, 0xd8, 0xcb, 0x56, //}
		0x78, 0x1e, 0xf2, 0xb5, //}  HMAC-SHA1 fingerprint
		0xb2, 0xd3, 0xf2, 0x49, //}
		0xc1, 0xb5, 0x71, 0xa2, //}
		0x80, 0x28, 0x00, 0x04, //   FINGERPRINT attribute header
		0xe5, 0x7a, 0x3b, 0xcf, //   CRC32 fingerprint
	}

	RFC5769SampleRequestMessage = MakeMessage(
		MessageType(0x0001),
		TransactionID{
			0xb7e7a701,
			0xbc34d686,
			0xfa87dfae,
		},
		[]Attribute{
			SoftwareAttribute{
				"STUN test client",
			},
			PriorityAttribute{
				0x6e0001ff,
			},
			IceControlledAttribute{
				0x932ff9b151263b36,
			},
			UsernameAttribute{
				"evtj:h6vY",
			},
			MessageIntegrityAttribute{
				HMAC: []byte{
					0x9a, 0xea, 0xa7, 0x0c,
					0xbf, 0xd8, 0xcb, 0x56,
					0x78, 0x1e, 0xf2, 0xb5,
					0xb2, 0xd3, 0xf2, 0x49,
					0xc1, 0xb5, 0x71, 0xa2,
				},
				raw: []byte{
					0x00, 0x01, 0x00, 0x58, //   Request type and message length
					0x21, 0x12, 0xa4, 0x42, //   Magic cookie
					0xb7, 0xe7, 0xa7, 0x01, //}
					0xbc, 0x34, 0xd6, 0x86, //}  Transaction ID
					0xfa, 0x87, 0xdf, 0xae, //}
					0x80, 0x22, 0x00, 0x10, //   SOFTWARE attribute header
					0x53, 0x54, 0x55, 0x4e, //}
					0x20, 0x74, 0x65, 0x73, //}  User-agent...
					0x74, 0x20, 0x63, 0x6c, //}  ...name
					0x69, 0x65, 0x6e, 0x74, //}
					0x00, 0x24, 0x00, 0x04, //   PRIORITY attribute header
					0x6e, 0x00, 0x01, 0xff, //   ICE priority value
					0x80, 0x29, 0x00, 0x08, //   ICE-CONTROLLED attribute header
					0x93, 0x2f, 0xf9, 0xb1, //}  Pseudo-random tie breaker...
					0x51, 0x26, 0x3b, 0x36, //}   ...for ICE control
					0x00, 0x06, 0x00, 0x09, //   USERNAME attribute header
					0x65, 0x76, 0x74, 0x6a, //}
					0x3a, 0x68, 0x36, 0x76, //}  Username (9 bytes) and padding (3 bytes)
					0x59, 0x20, 0x20, 0x20, //}
				},
			},
			FingerprintAttribute{
				0xe57a3bcf,
			},
		},
	)

	RFC5769SampleRequestWithLongTermCreds = []byte{
		0x00, 0x01, 0x00, 0x60, //    Request type and message length
		0x21, 0x12, 0xa4, 0x42, //    Magic cookie
		0x78, 0xad, 0x34, 0x33, // }
		0xc6, 0xad, 0x72, 0xc0, // }  Transaction ID
		0x29, 0xda, 0x41, 0x2e, // }
		0x00, 0x06, 0x00, 0x12, //    USERNAME attribute header
		0xe3, 0x83, 0x9e, 0xe3, // }
		0x83, 0x88, 0xe3, 0x83, // }
		0xaa, 0xe3, 0x83, 0x83, // }  Username value (18 bytes) and padding (2 bytes)
		0xe3, 0x82, 0xaf, 0xe3, // }
		0x82, 0xb9, 0x00, 0x00, // }
		0x00, 0x15, 0x00, 0x1c, //    NONCE attribute header
		0x66, 0x2f, 0x2f, 0x34, // }
		0x39, 0x39, 0x6b, 0x39, // }
		0x35, 0x34, 0x64, 0x36, // }
		0x4f, 0x4c, 0x33, 0x34, // }  Nonce value
		0x6f, 0x4c, 0x39, 0x46, // }
		0x53, 0x54, 0x76, 0x79, // }
		0x36, 0x34, 0x73, 0x41, // }
		0x00, 0x14, 0x00, 0x0b, //    REALM attribute header
		0x65, 0x78, 0x61, 0x6d, // }
		0x70, 0x6c, 0x65, 0x2e, // }  Realm value (11 bytes) and padding (1 byte)
		0x6f, 0x72, 0x67, 0x00, // }
		0x00, 0x08, 0x00, 0x14, //    MESSAGE-INTEGRITY attribute header
		0xf6, 0x70, 0x24, 0x65, // }
		0x6d, 0xd6, 0x4a, 0x3e, // }
		0x02, 0xb8, 0xe0, 0x71, // }  HMAC-SHA1 fingerprint
		0x2e, 0x85, 0xc9, 0xa2, // }
		0x8c, 0xa8, 0x96, 0x66, // }
	}

	RFC5769SampleRequestWithLongTermCredsMessage = MakeMessage(
		MessageType(0x0001),
		TransactionID{
			0x78ad3433,
			0xc6ad72c0,
			0x29da412e,
		},
		[]Attribute{
			UsernameAttribute{
				"\u30de\u30c8\u30ea\u30c3\u30af\u30b9",
			},
			NonceAttribute{
				"f//499k954d6OL34oL9FSTvy64sA",
			},
			RealmAttribute{
				"example.org",
			},
			MessageIntegrityAttribute{
				HMAC: []byte{
					0xf6, 0x70, 0x24, 0x65,
					0x6d, 0xd6, 0x4a, 0x3e,
					0x02, 0xb8, 0xe0, 0x71,
					0x2e, 0x85, 0xc9, 0xa2,
					0x8c, 0xa8, 0x96, 0x66,
				},
				raw: []byte{
					0x00, 0x01, 0x00, 0x60, //    Request type and message length
					0x21, 0x12, 0xa4, 0x42, //    Magic cookie
					0x78, 0xad, 0x34, 0x33, // }
					0xc6, 0xad, 0x72, 0xc0, // }  Transaction ID
					0x29, 0xda, 0x41, 0x2e, // }
					0x00, 0x06, 0x00, 0x12, //    USERNAME attribute header
					0xe3, 0x83, 0x9e, 0xe3, // }
					0x83, 0x88, 0xe3, 0x83, // }
					0xaa, 0xe3, 0x83, 0x83, // }  Username value (18 bytes) and padding (2 bytes)
					0xe3, 0x82, 0xaf, 0xe3, // }
					0x82, 0xb9, 0x00, 0x00, // }
					0x00, 0x15, 0x00, 0x1c, //    NONCE attribute header
					0x66, 0x2f, 0x2f, 0x34, // }
					0x39, 0x39, 0x6b, 0x39, // }
					0x35, 0x34, 0x64, 0x36, // }
					0x4f, 0x4c, 0x33, 0x34, // }  Nonce value
					0x6f, 0x4c, 0x39, 0x46, // }
					0x53, 0x54, 0x76, 0x79, // }
					0x36, 0x34, 0x73, 0x41, // }
					0x00, 0x14, 0x00, 0x0b, //    REALM attribute header
					0x65, 0x78, 0x61, 0x6d, // }
					0x70, 0x6c, 0x65, 0x2e, // }  Realm value (11 bytes) and padding (1 byte)
					0x6f, 0x72, 0x67, 0x00, // }
				}},
		},
	)
)
