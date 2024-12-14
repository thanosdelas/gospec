package spec

import (
  "testing"
)

func TestMain(testing *testing.T) {
  var spec = Spec{Testing: testing}

  spec.Expect(100).ToEq(100)
  spec.Expect(100).ToNotEq(200)
  spec.Expect("100").ToEq("100")
  spec.Expect("100").ToNotEq("200")
}
