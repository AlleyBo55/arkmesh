package erasure

import "errors"

type matrix [][]byte

func newMatrix(rows, columns int) matrix {
	built := make(matrix, rows)
	for row := range built {
		built[row] = make([]byte, columns)
	}
	return built
}

func identityMatrix(size int) matrix {
	built := newMatrix(size, size)
	for index := 0; index < size; index++ {
		built[index][index] = 1
	}
	return built
}

// vandermondeMatrix builds a rows x columns matrix whose element (r, c) is r^c.
// Every square submatrix of a Vandermonde matrix with distinct rows is
// invertible, which is what guarantees that any `columns` surviving shards can
// reconstruct the original data.
func vandermondeMatrix(rows, columns int) matrix {
	built := newMatrix(rows, columns)
	for row := 0; row < rows; row++ {
		for column := 0; column < columns; column++ {
			built[row][column] = exponent(byte(row), column)
		}
	}
	return built
}

func (left matrix) multiply(right matrix) matrix {
	product := newMatrix(len(left), len(right[0]))
	for row := range left {
		for column := range right[0] {
			var total byte
			for index := range right {
				total ^= mul(left[row][index], right[index][column])
			}
			product[row][column] = total
		}
	}
	return product
}

func (source matrix) subMatrix(rows []int) matrix {
	built := make(matrix, 0, len(rows))
	for _, row := range rows {
		built = append(built, append([]byte(nil), source[row]...))
	}
	return built
}

// invert returns the inverse of a square matrix using Gauss-Jordan elimination
// on the matrix augmented with the identity.
func (source matrix) invert() (matrix, error) {
	size := len(source)
	work := make(matrix, size)
	result := identityMatrix(size)
	for row := range source {
		if len(source[row]) != size {
			return nil, errors.New("erasure: matrix is not square")
		}
		work[row] = append([]byte(nil), source[row]...)
	}

	for column := 0; column < size; column++ {
		if work[column][column] == 0 {
			swapped := false
			for candidate := column + 1; candidate < size; candidate++ {
				if work[candidate][column] != 0 {
					work[column], work[candidate] = work[candidate], work[column]
					result[column], result[candidate] = result[candidate], result[column]
					swapped = true
					break
				}
			}
			if !swapped {
				return nil, errors.New("erasure: matrix is singular")
			}
		}
		pivot := work[column][column]
		if pivot != 1 {
			scale := inverse(pivot)
			scaleRow(work[column], scale)
			scaleRow(result[column], scale)
		}
		for row := 0; row < size; row++ {
			if row == column || work[row][column] == 0 {
				continue
			}
			factor := work[row][column]
			addScaledRow(work[row], work[column], factor)
			addScaledRow(result[row], result[column], factor)
		}
	}
	return result, nil
}

func scaleRow(row []byte, scale byte) {
	for index := range row {
		row[index] = mul(row[index], scale)
	}
}

func addScaledRow(target, source []byte, factor byte) {
	for index := range target {
		target[index] ^= mul(source[index], factor)
	}
}
