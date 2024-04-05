package locale

type Error struct {
	Page_Auth__Message_Success string ` default:"register was successful, check your email please!"`

	Mail struct {
		From_is_missing    string ` default:"mail 'from' is missing"`
		To_is_missing      string ` default:"mail 'to' is missing"`
		Subject_is_missing string ` default:"mail 'subject' is missing"`

		Send                     string ` default:"mail could not sent"`
		Add_attachment_on_attach string ` default:"could not add mail attachment"`
		Add_attachment_read_file string ` default:"could not add mail attachment, file not found or broken"`
	}
}
