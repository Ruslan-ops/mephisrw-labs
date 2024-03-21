package service

import (
	"backend/pkg/repository"
	"context"
	"math"
)

const Lab1AId = 1

type lab1aService struct {
	repo *repository.Repo
}

func NewLab1aService(repo *repository.Repo) *lab1aService {
	return &lab1aService{repo: repo}
}

func (s *lab1aService) CheckLab1AImportanceMatrix(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error) {
}

func (s *lab1aService) CheckLab1AImportanceMatrixFirstCriteria(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error) {
}

func (s *lab1aService) CheckLab1AImportanceMatrixSecondCriteria(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error) {
}

func (s *lab1aService) CheckLab1AImportanceMatrixThirdCriteria(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error) {
}

func (s *lab1aService) CheckLab1AImportanceMatrixFourthCriteria(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error) {
}

func (s *lab1aService) CheckLab1AChosenAlternative(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error) {
}

func (s *lab1aService) roundTwoDigits(value float64) float64 {
	return math.Round(value*100) / 100
}

func (s *lab1aService) calculatePriorityElement(arr []float64) float64 {
	product := 1.0
	for _, value := range arr {
		product *= value
	}
	geometricMean := math.Pow(product, 1/float64(len(arr)))
	return s.roundTwoDigits(geometricMean)
}

func (s *lab1aService) calculatePriorityVector(matrix [][]float64) []float64 {
	priorityVector := make([]float64, len(matrix))
	for i, arr := range matrix {
		priorityVector[i] = s.calculatePriorityElement(arr)
	}
	return priorityVector
}

func (s *lab1aService) calculateWeightElement(vectorElement float64, vectorSum float64) float64 {
	return s.roundTwoDigits(vectorElement / vectorSum)
}

func (s *lab1aService) calculateWeightVector(vector []float64, priorityVectorSum float64) []float64 {
	weightVector := make([]float64, len(vector))
	for i, element := range vector {
		weightVector[i] = s.calculateWeightElement(element, priorityVectorSum)
	}
	return weightVector
}

func (s *lab1aService) calculateMatrixWeightElement(matrixRow []float64, weightVector []float64) float64 {
	sum := 0.0
	for i, value := range matrixRow {
		sum += value * weightVector[i]
	}
	return s.roundTwoDigits(sum)
}

func (s *lab1aService) calculateMatrixWeightVector(matrix [][]float64, weightVector []float64) []float64 {
	matrixWeightVector := make([]float64, len(matrix))
	for i, row := range matrix {
		matrixWeightVector[i] = s.calculateMatrixWeightElement(row, weightVector)
	}
	return matrixWeightVector
}

func (s *lab1aService) calculateLambdaVector(weightVector []float64, matrixWeightVector []float64) []float64 {
	lambdaVector := make([]float64, len(weightVector))
	for i, value := range matrixWeightVector {
		lambdaVector[i] = s.roundTwoDigits(value / weightVector[i])
	}
	return lambdaVector
}

func (s *lab1aService) calculateEigenvalue(lambdaVector []float64) float64 {
	max := lambdaVector[0]
	for _, value := range lambdaVector {
		if value > max {
			max = value
		}
	}
	return max
}

func (s *lab1aService) calculateConsistencyIndex(eigenvalue float64, lambdaVector []float64) float64 {
	return s.roundTwoDigits(eigenvalue-float64(len(lambdaVector))) / (float64(len(lambdaVector) - 1))
}

func (s *lab1aService) calculateConsistencyRatio(matrix [][]float64, consistencyIndex float64) float64 {
	return s.roundTwoDigits(consistencyIndex / s.mathExpectationsOfConsistencyIndex[len(matrix)])
}

func (s *lab1aService) isConsistencyAcceptable(consistencyRatio float64) bool {
	return math.Abs(consistencyRatio) < 0.1
}

func (s *lab1aService) calculateWeightsRatioMatrix(weightVector []float64) [][]float64 {
	matrix := make([][]float64, len(weightVector))
	for _, weight := range weightVector {
		row := make([]float64, len(weightVector))
		for i, value := range weightVector {
			row[i] = s.roundTwoDigits(weight / value)
		}
		matrix = append(matrix, row)
	}
	return matrix
}

func (s *lab1aService) calculateCorrectedWeightRatioMatrix(matrix [][]float64, weightsRatioMatrix [][]float64) [][]float64 {
	newMatrix := make([][]float64, len(matrix))
	for i, row := range matrix {
		newRow := make([]float64, len(row))
		for j, value := range row {
			newRow[j] = math.Abs(s.roundTwoDigits(value - weightsRatioMatrix[i][j]))
		}
		newMatrix = append(newMatrix, newRow)
	}
	return newMatrix
}

func (s *lab1aService) findMaxSumIndex(weightRatioMatrix [][]float64) int {
	maxSum := 0.0
	maxSumIndex := -1
	for i, row := range weightRatioMatrix {
		sum := 0.0
		for _, value := range row {
			sum += value
		}
		if sum > maxSum {
			maxSum = sum
			maxSumIndex = i
		}
	}
	return maxSumIndex
}

func (s *lab1aService) calculateNewCorrectedMatrix(matrix [][]float64, weightRatioMatrix [][]float64, indexToChange int) [][]float64 {
	newMatrix := make([][]float64, len(matrix))
	for i, row := range matrix {
		if i == indexToChange {
			newMatrix = append(newMatrix, weightRatioMatrix[i])
		} else {
			newRow := make([]float64, len(row))
			for j, value := range row {
				if j == indexToChange {
					newRow[j] = s.roundTwoDigits(1 / weightRatioMatrix[indexToChange][i])
				} else {
					newRow[j] = value
				}
			}
			newMatrix = append(newMatrix, newRow)
		}
	}
	return newMatrix
}

func (s *lab1aService) solveMatrix(matrix [][]float64) {
	priorityVector := s.calculatePriorityVector(matrix)
	priorityVectorSum := 0.0
	for _, value := range priorityVector {
		priorityVectorSum += value
	}
	weightVector := s.calculateWeightVector(priorityVector, priorityVectorSum)
	matrixWeightVector := s.calculateMatrixWeightVector(matrix, weightVector)
	lambdaVector := s.calculateLambdaVector(weightVector, matrixWeightVector)
	eigenvalue := s.calculateEigenvalue(lambdaVector)
	consistencyIndex := s.calculateConsistencyIndex(eigenvalue, lambdaVector)
	consistencyRatio := s.calculateConsistencyRatio(matrix, consistencyIndex)

	if !s.isConsistencyAcceptable(consistencyRatio) {
		weightsRatioMatrix := s.calculateWeightsRatioMatrix(weightVector)
		weightsCorrectedRatioMatrix := s.calculateCorrectedWeightRatioMatrix(matrix, weightsRatioMatrix)
		maxSumIndex := s.findMaxSumIndex(weightsCorrectedRatioMatrix)
		newVariant := s.calculateNewCorrectedMatrix(matrix, weightsRatioMatrix, maxSumIndex)
		s.solveMatrix(newVariant)
	}
}
