package shader

var Chunk_Shader = `//kage:unit pixels
			package main

			var R float 
			var G float
			var B float

			var RR float
			var GG float
			var BB float

			func Fragment(targetCoords vec4, srcPos vec2, _ vec4) vec4 {
				col := imageSrc0At(srcPos.xy)
				if col.x >= 163.0/255 &&col.x <= 165.0/255{
					return vec4(R, G, B, 255)
				}
				if col.x >= 132.0/255 && col.x <= 134.0/255{
					return vec4(RR, GG, BB, 255)
				}

				if col.x >= 126.0/255 && col.x <= 128.0/255 {
					var retun_col vec4 

					retun_col.r = ((R+RR)/1.5) 
					retun_col.g = ((G+GG)/1.5) 
					retun_col.b = ((B+BB)/1.5) 
	
					retun_col.r /= 2
					retun_col.g /= 2
					retun_col.b /= 2
					retun_col.a = 0.5

					return retun_col
				}

				return col
			}
`
