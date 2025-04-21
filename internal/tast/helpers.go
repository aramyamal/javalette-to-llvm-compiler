package tast

import "slices"

// GuaranteesReturn checks if the given statement guarantees a return on all
// paths when traversing all children nodes of the statement node TAST.
func GuaranteesReturn(stm Stm) bool {
	switch s := stm.(type) {
	case *ReturnStm:
		return true
	case *BlockStm:
		// a block guarantees return if at least one statement guarantees return
		return slices.ContainsFunc(s.Stms, GuaranteesReturn)
	case *IfStm:
		// if statement guarantees return only if both branches guarantee return
		if s.ElseStm == nil {
			return false // no else branch means no guarantee
		}
		return GuaranteesReturn(s.ThenStm) && GuaranteesReturn(s.ElseStm)
	default:
		return false
	}
}
