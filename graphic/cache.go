package graphic

// Invalidate the cache of this object
// Implying that the cache of the parents objects all have to be removed too, because they would presumablly change
func InvalidateRenderCache(object Object) {
	// if already invalidated then return
	if !object.InvalidateRenderCache() {
		return
	}

	// otherwise invalidate all the objects that need this object too
	asChild, isChildObject := object.(ChildObject)
	if isChildObject {
		parent := asChild.GetParent()
		if parent != nil {
			InvalidateRenderCache(parent)
		}
	}
}

// Invalidate the render cache of this object but also the objects that make up this objects and so on
func InvalidateRenderCacheRecursive(object Object) {
	InvalidateRenderCache(object)

	_, isComposite := object.(CompositeObject)
	if isComposite {
		object.(CompositeObject).ForEach(func(child Object) {
			InvalidateRenderCacheRecursive(child)
		})
	}
}
