// +build ignore

package wishlist

import (
	"github.com/corestoreio/csfw/config/element"
	"github.com/corestoreio/csfw/config/model"
)

// Path will be initialized in the init() function together with ConfigStructure.
var Path *PkgPath

// PkgPath global configuration struct containing paths and how to retrieve
// their values and options.
type PkgPath struct {
	model.PkgPath
	// WishlistEmailEmailIdentity => Email Sender.
	// Path: wishlist/email/email_identity
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Identity
	WishlistEmailEmailIdentity model.Str

	// WishlistEmailEmailTemplate => Email Template.
	// Email template chosen based on theme fallback when "Default" option is
	// selected.
	// Path: wishlist/email/email_template
	// SourceModel: Otnegam\Config\Model\Config\Source\Email\Template
	WishlistEmailEmailTemplate model.Str

	// WishlistEmailNumberLimit => Max Emails Allowed to be Sent.
	// 10 by default. Max - 10000
	// Path: wishlist/email/number_limit
	WishlistEmailNumberLimit model.Str

	// WishlistEmailTextLimit => Email Text Length Limit.
	// 255 by default
	// Path: wishlist/email/text_limit
	WishlistEmailTextLimit model.Str

	// WishlistGeneralActive => Enabled.
	// Path: wishlist/general/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Yesno
	WishlistGeneralActive model.Bool

	// WishlistWishlistLinkUseQty => Display Wish List Summary.
	// Path: wishlist/wishlist_link/use_qty
	// SourceModel: Otnegam\Wishlist\Model\Config\Source\Summary
	WishlistWishlistLinkUseQty model.Str

	// RssWishlistActive => Enable RSS.
	// Path: rss/wishlist/active
	// SourceModel: Otnegam\Config\Model\Config\Source\Enabledisable
	RssWishlistActive model.Bool
}

// NewPath initializes the global Path variable. See init()
func NewPath(cfgStruct element.SectionSlice) *PkgPath {
	return (&PkgPath{}).init(cfgStruct)
}

func (pp *PkgPath) init(cfgStruct element.SectionSlice) *PkgPath {
	pp.Lock()
	defer pp.Unlock()
	pp.WishlistEmailEmailIdentity = model.NewStr(`wishlist/email/email_identity`, model.WithConfigStructure(cfgStruct))
	pp.WishlistEmailEmailTemplate = model.NewStr(`wishlist/email/email_template`, model.WithConfigStructure(cfgStruct))
	pp.WishlistEmailNumberLimit = model.NewStr(`wishlist/email/number_limit`, model.WithConfigStructure(cfgStruct))
	pp.WishlistEmailTextLimit = model.NewStr(`wishlist/email/text_limit`, model.WithConfigStructure(cfgStruct))
	pp.WishlistGeneralActive = model.NewBool(`wishlist/general/active`, model.WithConfigStructure(cfgStruct))
	pp.WishlistWishlistLinkUseQty = model.NewStr(`wishlist/wishlist_link/use_qty`, model.WithConfigStructure(cfgStruct))
	pp.RssWishlistActive = model.NewBool(`rss/wishlist/active`, model.WithConfigStructure(cfgStruct))

	return pp
}
