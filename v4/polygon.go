package sameriver

import (
	"math"
)

// Polygon represents a 2D polygon defined by a list of vertices.
type Polygon struct {
	Vertices []Vec2D
}

// DistanceToSide calculates the shortest distance from a point to any side of the polygon.
func (p *Polygon) DistanceToSide(point Vec2D) float64 {
	minDistance := math.MaxFloat64

	for i := 0; i < len(p.Vertices); i++ {
		// Calculate distance to line segment formed by vertices[i] and vertices[(i+1)%n]
		d := pointToLineSegmentDistance(point, p.Vertices[i], p.Vertices[(i+1)%len(p.Vertices)])

		// Update minimum distance if the calculated distance is smaller
		if d < minDistance {
			minDistance = d
		}
	}

	return minDistance
}

// DistanceToVertex calculates the distance from a point to the closest vertex of the polygon.
func (p *Polygon) DistanceToVertex(point Vec2D) float64 {
	minDistance := math.MaxFloat64

	for _, vertex := range p.Vertices {
		// Calculate distance to vertex
		_, _, d := point.Distance(vertex)

		// Update minimum distance if the calculated distance is smaller
		if d < minDistance {
			minDistance = d
		}
	}

	return minDistance
}

// pointToLineSegmentDistance calculates the distance from a point to a line segment formed by two vertices.
func pointToLineSegmentDistance(point, lineStart, lineEnd Vec2D) float64 {
	// Calculate vector representing the line segment
	lineVector := Vec2D{lineEnd.X - lineStart.X, lineEnd.Y - lineStart.Y}

	// Calculate vector representing the point to lineStart
	pointToLineStart := Vec2D{point.X - lineStart.X, point.Y - lineStart.Y}

	// Calculate the projection of pointToLineStart onto the lineVector
	projection := (pointToLineStart.X*lineVector.X + pointToLineStart.Y*lineVector.Y) / (lineVector.X*lineVector.X + lineVector.Y*lineVector.Y)

	// Check if the projection is outside the line segment, and if so, calculate distance to the closest vertex
	if projection < 0 {
		_, _, d := point.Distance(lineStart)
		return d
	} else if projection > 1 {
		_, _, d := point.Distance(lineEnd)
		return d
	}

	// Calculate the closest point on the line segment to the given point
	closestPoint := Vec2D{lineStart.X + projection*lineVector.X, lineStart.Y + projection*lineVector.Y}

	// Calculate the distance between the given point and the closest point on the line segment
	_, _, d := point.Distance(closestPoint)
	return d
}
