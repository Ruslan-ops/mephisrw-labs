package service

import (
	"backend/pkg/repository"
	"math"
)

const Lab1AId = 1

type lab1aService struct {
	repo *repository.Repo
}

func NewLab1aService(repo *repository.Repo) *lab1aService {
	return &lab1aService{repo: repo}
}

func roundTwoDigits(value float64) float64 {
 return math.Round(value*100) / 100
}

func calculatePriorityElement(arr []float64) float64 {
 product := 1.0
 for _, value := range arr {
  product *= value
 }
 geometricMean := math.Pow(product, 1/float64(len(arr)))
 return roundTwoDigits(geometricMean)
}

func calculatePriorityVector(matrix [][]float64) []float64 {
 priorityVector := make([]float64, len(matrix))
 for i, arr := range matrix {
  priorityVector[i] = calculatePriorityElement(arr)
 }
 return priorityVector
}

func calculateWeightElement(vectorElement float64, vectorSum float64) float64 {
 return roundTwoDigits(vectorElement / vectorSum)
}

func calculateWeightVector(vector []float64, priorityVectorSum float64) []float64 {
 weightVector := make([]float64, len(vector))
 for i, element := range vector {
  weightVector[i] = calculateWeightElement(element, priorityVectorSum)
 }
 return weightVector
}

func calculateMatrixWeightElement(matrixRow []float64, weightVector []float64) float64 {
 sum := 0.0
 for i, value := range matrixRow {
  sum += value * weightVector[i]
 }
 return roundTwoDigits(sum)
}

func calculateMatrixWeightVector(matrix [][]float64, weightVector []float64) []float64 {
 matrixWeightVector := make([]float64, len(matrix))
 for i, row := range matrix {
  matrixWeightVector[i] = calculateMatrixWeightElement(row, weightVector)
 }
 return matrixWeightVector
}

func calculateLambdaVector(weightVector []float64, matrixWeightVector []float64) []float64 {
 lambdaVector := make([]float64, len(weightVector))
 for i, value := range matrixWeightVector {
  lambdaVector[i] = roundTwoDigits(value / weightVector[i])
 }
 return lambdaVector
}

func calculateEigenvalue(lambdaVector []float64) float64 {
 max := lambdaVector[0]
 for _, value := range lambdaVector {
  if value > max {
   max = value
  }
 }
 return max
}

func calculateConsistencyIndex(eigenvalue float64, lambdaVector []float64) float64 {
 return roundTwoDigits(eigenvalue - float64(len(lambdaVector))) / (float64(len(lambdaVector) - 1))
}

func calculateConsistencyRatio(matrix [][]float64, consistencyIndex float64) float64 {
 return roundTwoDigits(consistencyIndex / mathExpectationsOfConsistencyIndex[len(matrix)])
}

func isConsistencyAcceptable(consistencyRatio float64) bool {
 return math.Abs(consistencyRatio) < 0.1
}

func calculateWeightsRatioMatrix(weightVector []float64) [][]float64 {
 matrix := make([][]float64, len(weightVector))
 for _, weight := range weightVector {
  row := make([]float64, len(weightVector))
  for i, value := range weightVector {
   row[i] = roundTwoDigits(weight / value)
  }
  matrix = append(matrix, row)
 }
 return matrix
}

func calculateCorrectedWeightRatioMatrix(matrix [][]float64, weightsRatioMatrix [][]float64) [][]float64 {
 newMatrix := make([][]float64, len(matrix))
 for i, row := range matrix {
  newRow := make([]float64, len(row))
  for j, value := range row {
   newRow[j] = math.Abs(roundTwoDigits(value - weightsRatioMatrix[i][j]))
  }
  newMatrix = append(newMatrix, newRow)
 }
 return newMatrix
}

func findMaxSumIndex(weightRatioMatrix [][]float64) int {
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

func calculateNewCorrectedMatrix(matrix [][]float64, weightRatioMatrix [][]float64, indexToChange int) [][]float64 {
 newMatrix := make([][]float64, len(matrix))
 for i, row := range matrix {
  if i == indexToChange {
   newMatrix = append(newMatrix, weightRatioMatrix[i])
  } else {
   newRow := make([]float64, len(row))
   for j, value := range row {
    if j == indexToChange {
     newRow[j] = roundTwoDigits(1 / weightRatioMatrix[indexToChange][i])
    } else {
     newRow[j] = value
    }
   }
   newMatrix = append(newMatrix, newRow)
  }
 }
 return newMatrix
}

func solveMatrix(matrix [][]float64) {
 priorityVector := calculatePriorityVector(matrix)
 priorityVectorSum := 0.0
 for _, value := range priorityVector {
  priorityVectorSum += value
 }
 weightVector := calculateWeightVector(priorityVector, priorityVectorSum)
 matrixWeightVector := calculateMatrixWeightVector(matrix, weightVector)
 lambdaVector := calculateLambdaVector(weightVector, matrixWeightVector)
 eigenvalue := calculateEigenvalue(lambdaVector)
 consistencyIndex := calculateConsistencyIndex(eigenvalue, lambdaVector)
 consistencyRatio := calculateConsistencyRatio(matrix, consistencyIndex)

 if !isConsistencyAcceptable(consistencyRatio) {
  weightsRatioMatrix := calculateWeightsRatioMatrix(weightVector)
  weightsCorrectedRatioMatrix := calculateCorrectedWeightRatioMatrix(matrix, weightsRatioMatrix)
  maxSumIndex := findMaxSumIndex(weightsCorrectedRatioMatrix)
  newVariant := calculateNewCorrectedMatrix(matrix, weightsRatioMatrix, maxSumIndex)
  solveMatrix(newVariant)
 }
}