TODO
====

- [ ] Exit if cloud not found by name
- [ ] Take IPAM network name and subnet as input for creating VS
- [ ] Exit if IPAM network and subnet is not valid, doesn't belong to
  the cloud
- [ ] Multiple pools for multiple ports per container
- [ ] Take labels and configure VS accordingly
    - [ ] Enable SSL based on label
    - [ ] SSL Termination: certificate configuration as label
    - [ ] Other SSL related items
- [ ] Active monitoring set up based on protocol and port
- [ ] Periodic Sync of VS and Pool Members with Avi
- [ ] Use SSL Verify while connecting to Avi controller 
- [ ] Support for separate data and control networks
