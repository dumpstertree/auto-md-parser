package main

import (
	"encoding/json"
	"errors"
)

func (f *OrderedLayout) UnmarshalJSON(b []byte) error {
	type fleet OrderedLayout

	// guard - failed to unmarshal
	err := json.Unmarshal(b, (*fleet)(f))
	if err != nil {
		return err
	}

	// iterate over each entry in raw string
	for _, raw := range f.Sections {

		var v BaseSubsection

		// guard - failed to unmarshal
		err = json.Unmarshal(raw, &v)
		if err != nil {
			return err
		}

		// assign based on type
		var i ISubsection
		switch v.Type {
		case "paragraph":
			i = &Paragraph{}
		case "orderedlist":
			i = &OrderedList{}
		case "unorderedlist":
			i = &UnorderedList{}
		case "header":
			i = &Header{}
		case "subheader":
			i = &Subheader{}
		case "table":
			i = &Table{}
		case "quote":
			i = &Quote{}
		case "code":
			i = &Code{}
		case "subpage":
			i = &Subpage{}
		default:
			return errors.New("unknown type")
		}

		// guard - failed to assign
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}

		// add
		f.LayoutSubsection = append(f.LayoutSubsection, i)
	}
	return nil
}

func (f *Subpage) UnmarshalJSON(b []byte) error {
	type fleet Subpage

	// guard - failed to unmarshal
	err := json.Unmarshal(b, (*fleet)(f))
	if err != nil {
		return err
	}

	// iterate over each entry in raw string
	for _, raw := range f.Sections {

		var v BaseSubsection

		// guard - failed to unmarshal
		err = json.Unmarshal(raw, &v)
		if err != nil {
			return err
		}

		// assign based on type
		var i ISubSubsection
		switch v.Type {
		case "paragraph":
			i = &Paragraph{}
		// case "orderedlist":
		// 	i = &OrderedList{}
		// case "unorderedlist":
		// 	i = &UnorderedList{}
		case "header":
			i = &Header{}
		case "subheader":
			i = &Subheader{}
		// case "table":
		// 	i = &Table{}
		// case "quote":
		// 	i = &Quote{}
		// case "code":
		// 	i = &Code{}
		// case "subpage":
		// i = &Subpage{}
		default:
			return errors.New("unknown type")
		}

		// guard - failed to assign
		err = json.Unmarshal(raw, i)
		if err != nil {
			return err
		}

		// add
		f.LayoutSubsection = append(f.LayoutSubsection, i)
	}
	return nil
}
