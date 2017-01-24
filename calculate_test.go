package calculate

import (
	"testing"

	"fmt"
	"github.com/thisisfineio/common"
	"github.com/alistanis/awstools/awsregions"
	. "github.com/smartystreets/goconvey/convey"
)

/*
 * The following tests will hit actual AWS Endpoints, so in order to run them you must actually have an account and
 * the proper configurations set up
 */
func TestNewEC2(t *testing.T) {
	Convey("We can get a new EC2 struct", t, func() {
		ec2 := NewEC2(awsregions.USEast1)
		So(ec2, ShouldNotBeNil)
	})
}

func TestEC2_SetRegion(t *testing.T) {
	Convey("We can change the region of an EC2 struct", t, func() {
		ec2 := NewEC2(awsregions.AsiaNortheast1)
		So(ec2, ShouldNotBeNil)
		So(*ec2.service.Config.Region, ShouldEqual, awsregions.AsiaNortheast1)
		ec2.SetRegion(awsregions.USEast1)
		So(*ec2.service.Config.Region, ShouldEqual, awsregions.USEast1)
	})
}

func TestEC2_DescribeInstances(t *testing.T) {
	Convey("We can describe instances in our account", t, func() {
		ec2 := NewEC2(awsregions.USEast1)
		So(ec2, ShouldNotBeNil)
	})
}

func TestNewCompute(t *testing.T) {
	Convey("We can get a new compute interface with a valid provider, and an error when an invalid one is given", t, func() {
		c, err := NewCompute(common.AwsProvider, awsregions.USEast1)
		So(err, ShouldBeNil)
		So(c.Provider(), ShouldEqual, common.AwsProvider)
		_, err = NewCompute(-1, "no region")
		e := err.(*common.Error)
		So(e.Code, ShouldEqual, common.InvalidProviderErrCode)

		Convey("We can test provider functions", func() {
			instances, err := c.DescribeInstances(nil)
			So(err, ShouldBeNil)
			fmt.Println(instances)
		})
	})
}
