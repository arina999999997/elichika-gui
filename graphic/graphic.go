package graphic

// This package is an abstraction layer of various graphic concepts, just so we can switch to a different graphic library if necessary
// We follow the conventions:
// - Coordinate x refer to left - right, the higher x is, the more to right a point is.
// - Coordinate y refer to up - down, the higher y is, the more to the bottom a point is.
// - When refering to size, the width (x wise) come before the height (y wise).
// - Things that need to be drawn are refered to as objects.
// - Each object is treated as a rectangle. and has an internal resolution of HEIGHT * WIDTH.
// - There is a special root object that is the whole window display.
// - Each object can be a texture as is:
//   - These can be refered to as prime objects.
//   - The texture can be loaded from an image, or is a solid color, or even transparency.
//   - If loaded from image, the texture and its internal resolution will be that image.
//   - If filled with solid color, then the texture will require HEIGHT and WIDTH to be provided.
//   - When draw, a texture is just drawn as is, on top of the canvas.
// - If an object isn't a texture, then it must contain children objects:
//   - These can be refered to as composite objects.
//   - The children are kept in a list, based on the order of drawing:
// 	   - To draw a composite objects, draw the childrens from first to last on the list.
//   - Each child can have its own internal resolution, however:
//     - It must be anchored to the parent based on the parent's internal resolution.
//     - When anchoring, the top left corner along with the size in term of the parent's resolution must be provided.
//     - The children's anchoring point can be anywhere, even outside the range of the parent.
//     - The children's size can be anything, even higher than the size of the parent.
//   - The childrens can in turn have children too, forming a tree structure.
// - (TODO) Other than the above, it's possible to add transformation object:
//   - These objects would take another object and transform it.
//   - For example, mirroring, or stretch / squeezing.
// - In any cases, all the object must form a strict "tree structure":
//   - Prime objects can be reused and shared.
//   - When a prime object is anchored to a parent object, the parent object gets its own "copy" of the prime object.
//   - So every instance of the prime object being used is different.
//   - To modify a specific instance of a prime object, we need to modify the actual instance, not the original object before anchoring.
//   - In practice, we shouldn't be modifying prime objects when running:
//     - It's obviously the case if this is a texture loaded from disk.
//     - If the texture is uploaded by the user, we should be replacing the object with another object, not reloading.
//     - The same for reloading / regenerating texture from code.
//
import (
	"runtime"
)

func init() { runtime.LockOSThread() }

// TODO(memory): This thing have lots of memory leak in it, so don't use it for anything that can result in leaks beyond controls
