# Above the Clouds Review
## David Cleaver

Above the Clouds describes the cloud as *Utility Computing*. The paper outlines
that the utility nature of raw compute is a necessary condition for the rise of
the cloud. Much like turning on the power switch to get electricity into your
home, the cloud user needs to be able to access more cpu, storage, and network as
soon as they need it. In the same way that one is billed for power on the amount
of power being used, the cloud bills based on usage. The paper reflects this by
calling out the necessary pattern of scaling up and down services in the cloud
to better manage costs. In the same way that you would turn the lights off in a room
when you aren't using them to save on electricity costs. The scale at which
cloud datacenters can be and are built makes these resources appear effectively
infinite to the average cloud user. When spinning up new cloud usage, there is
no need to worry that Amazon, Google, or Microsoft will run out of capacity.

A "Lift and Shift" migration to the cloud allows a software engineer to treat
the cloud as a utility like power or water. Essentially, this method of
migration involves moving the application as is into the cloud without making
any changes to how it operates.
[[1](https://www.cloudzero.com/blog/lift-and-shift/)] This treats the cloud just
like a new datacenter and is very similar to simply leaving the water running,
or not turning out lights when you aren't using them. You will spend more money,
but the benefits can include moving faster to the cloud or getting familiar with
the cloud.

Cost optimizations in the cloud can lead software engineers to make odd
architectural decisions to support lower cost cloud options. For example, it costs 
money to transfer data between AWS availibility zones, but S3 can be used to
transfer that data for a much lower monthly cost.
[[2](https://www.bitsand.cloud/posts/slashing-data-transfer-costs/)] The
software engineer that makes these kinds of decisions is not able to consider
the cloud as a pure utility that they can utilize. Instead they are forced to
consider ways to try to work around the provider's billing by using their
services creatively rather than naturally. These choices also make the
architecture very specific to the cloud provider and difficult to move to
another equivalent provider.
