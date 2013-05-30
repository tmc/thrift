/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package thrift

func compareInt(i, j int) int {
	if i > j {
		return 1
	}
	if i < j {
		return -1
	}
	return 0
}

func compareInt16(i, j int16) int {
	if i > j {
		return 1
	}
	if i < j {
		return -1
	}
	return 0
}

func compareInt32(i, j int32) int {
	if i > j {
		return 1
	}
	if i < j {
		return -1
	}
	return 0
}

func compareInt64(i, j int32) int {
	if i > j {
		return 1
	}
	if i < j {
		return -1
	}
	return 0
}

func compareStringArray(i, j []string) int {
	if cmp := compareInt(len(i), len(j)); cmp != 0 {
		return cmp
	}
	size := len(i)
	for k := 0; k < size; k++ {
		if cmp := compareString(i[k], j[k]); cmp != 0 {
			return cmp
		}
	}
	return 0
}

func compareString(i, j string) int {
	if i > j {
		return 1
	}
	if i < j {
		return -1
	}
	return 0
}

func compareFloat(i, j float32) int {
	if i > j {
		return 1
	}
	if i < j {
		return -1
	}
	return 0
}

func compareDouble(i, j float64) int {
	if i > j {
		return 1
	}
	if i < j {
		return -1
	}
	return 0
}

func compareByte(i, j byte) int {
	if i > j {
		return 1
	}
	if i < j {
		return -1
	}
	return 0
}

func compareBool(i, j bool) int {
	if i {
		if j {
			return 0
		}
		return 1
	}
	if j {
		return -1
	}
	return 0
}
