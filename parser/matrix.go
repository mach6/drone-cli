package parser

import (
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	limitTags = 10
	limitAxis = 25
)

// Matrix represents the build matrix.
type Matrix map[string][]string

// Axis represents a single permutation of entries
// from the build matrix.
type Axis map[string]string

// String returns a string representation of an
// Axis as a comma-separated list of environment
// variables.
func (a Axis) String() string {
	var envs []string
	for k, v := range a {
		envs = append(envs, k+"="+v)
	}
	return strings.Join(envs, ",")
}

// ParseMatrix parses the Matrix section of the yaml
// file and returns a list of axis.
func ParseMatrix(raw string) ([]Axis, error) {
	matrix, err := parseMatrix(raw)
	if err != nil {
		return nil, err
	}

	// if not a matrix build return an array
	// with just the single axis.
	if len(matrix) == 0 {
		return nil, nil
	}

	return CalcMatrix(matrix), nil
}

// CalcMatrix calculates the permutations for the
// build matrix.
//
// Note that this method will cap the number of
// permutations to 25 to prevent an overly
// expensive calculation.
func CalcMatrix(matrix Matrix) []Axis {
	// calculate number of permutations and
	// extract the list of tags
	// (ie go_version, redis_version, etc)
	var perm int
	var tags []string
	for k, v := range matrix {
		perm *= len(v)
		if perm == 0 {
			perm = len(v)
		}
		tags = append(tags, k)
	}

	// structure to hold the transformed
	// result set
	axisList := []Axis{}

	// for each axis calculate the uniqe
	// set of values that should be used.
	for p := 0; p < perm; p++ {
		axis := map[string]string{}
		decr := perm
		for i, tag := range tags {
			elems := matrix[tag]
			decr = decr / len(elems)
			elem := p / decr % len(elems)
			axis[tag] = elems[elem]

			// enforce a maximum number of tags
			// in the build matrix.
			if i > limitTags {
				break
			}
		}

		// append to the list of axis.
		axisList = append(axisList, axis)

		// enforce a maximum number of axis
		// that should be calculated.
		if p > limitAxis {
			break
		}
	}

	return axisList
}

// helper function to parse the Matrix data from
// the raw yaml file.
func parseMatrix(raw string) (Matrix, error) {
	data := struct {
		Matrix map[string][]string
	}{}
	err := yaml.Unmarshal([]byte(raw), &data)
	return data.Matrix, err
}
