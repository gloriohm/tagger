package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"tagger/internal/compare"
	"tagger/internal/mbfetcher"
	"tagger/internal/workflow"
)

type UserCapture struct {
	in  *bufio.Reader
	out io.Writer
}

func NewUserCapture(in io.Reader, out io.Writer) *UserCapture {
	return &UserCapture{
		in:  bufio.NewReader(in),
		out: out,
	}
}

func (s *UserCapture) SelectRelease(releases []mbfetcher.ReleaseInfo) (workflow.ReleaseSelection, error) {
	for i, release := range releases {
		fmt.Fprintf(s.out, "%d. %s\n", i+1, release.Quality)
	}

	fmt.Fprintln(s.out, "n. Enter release ID manually")

	for {
		fmt.Fprint(s.out, "> ")

		input, err := s.in.ReadString('\n')
		if err != nil {
			return workflow.ReleaseSelection{}, err
		}

		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "n") {
			fmt.Fprint(s.out, "Release ID: ")

			id, err := s.in.ReadString('\n')
			if err != nil {
				return workflow.ReleaseSelection{}, err
			}

			id = strings.TrimSpace(id)
			if id == "" {
				fmt.Fprintln(s.out, "Release ID cannot be empty")
				continue
			}

			return workflow.ReleaseSelection{
				ReleaseID: id,
			}, nil
		}
		index, err := strconv.Atoi(input)
		if err != nil || index < 1 || index > len(releases) {
			fmt.Fprintln(s.out, "Invalid selection")
			continue
		}

		return workflow.ReleaseSelection{
			ReleaseID: releases[index-1].ID,
		}, nil
	}
}

func (c *UserCapture) ReviewRelease(comp compare.Comparision) (workflow.Decision, error) {
	fmt.Fprintf(c.out, "Total distance: %d\n", comp.Distance)

	for _, diff := range comp.Differences {
		fmt.Fprintf(c.out, "field/distance: %s --- %d\n", diff.Field, diff.Distance)
	}

	fmt.Fprintln(c.out, "y/n --- approve/reject")

	for {
		fmt.Fprint(c.out, "> ")

		input, err := c.in.ReadString('\n')
		if err != nil {
			return workflow.Decision{}, err
		}

		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "n") {
			fmt.Fprint(c.out, "Enter new release ID: ")

			id, err := c.in.ReadString('\n')
			if err != nil {
				return workflow.Decision{}, err
			}

			id = strings.TrimSpace(id)
			if id == "" {
				fmt.Fprintln(c.out, "Release ID cannot be empty")
				continue
			}

			return workflow.Decision{
				Approved:  false,
				ReleaseID: id,
			}, nil
		} else if strings.EqualFold(input, "y") {
			return workflow.Decision{
				Approved: true,
			}, nil
		} else {
			fmt.Fprintln(c.out, "Invalid input")
			continue
		}
	}
}
