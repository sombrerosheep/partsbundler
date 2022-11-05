package main

type BundlerService struct {
	Parts IPartService
	Kits  IKitService
}